package forge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Forgejo implements the Forge interface for Forgejo/Gitea.
type Forgejo struct {
	baseURL string
	token   string
	client  *http.Client
}

// ForgejoConfig holds configuration for Forgejo forge.
type ForgejoConfig struct {
	BaseURL string
	Token   string
}

// NewForgejo creates a new Forgejo forge client.
func NewForgejo(baseURL, token string) *Forgejo {
	return &Forgejo{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: 30 * time.Second}, //nolint:mnd
	}
}

// Name returns the name of the forge.
func (f *Forgejo) Name() string {
	return "forgejo"
}

// forgejoGetPR is the API response for a pull request.
type forgejoPR struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	Merged    bool   `json:"merged"`
	MergeBase string `json:"merge_base"`
	MergedAt  string `json:"merged_at"`
	MergeSHA  string `json:"merge_commit_sha"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
	Head struct {
		SHA string `json:"sha"`
		Ref string `json:"ref"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// forgejoCommit is the API response for a commit.
type forgejoCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	Parents []struct {
		SHA string `json:"sha"`
	} `json:"parents"`
}

// forgejoError is the API error response.
type forgejoError struct {
	Message string `json:"message"`
}

// parseForgejoError extracts a clean error message from API response.
func parseForgejoError(body []byte) string {
	var errResp forgejoError
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Message != "" {
		return errResp.Message
	}
	// Fallback to raw body, but clean it up
	return strings.TrimSpace(string(body))
}

// GetPR retrieves information about a pull request by number.
func (f *Forgejo) GetPR(ctx context.Context, owner, repo string, number int) (*PRInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/pulls/%d", f.baseURL, owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if f.token != "" {
		req.Header.Set("Authorization", "token "+f.token)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR #%d: %w", number, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get PR #%d: %s (%s)", number, resp.Status, parseForgejoError(body))
	}

	var pr forgejoPR
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode PR response: %w", err)
	}

	if !pr.Merged {
		return nil, fmt.Errorf("PR #%d is not merged", number)
	}

	// Get merge commit to check if squashed.
	mergeCommit, err := f.GetCommit(ctx, owner, repo, pr.MergeSHA)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge commit: %w", err)
	}

	squashed := len(mergeCommit.Parents) == 1

	mergedAt, _ := time.Parse(time.RFC3339, pr.MergedAt)

	info := &PRInfo{
		Number:      pr.Number,
		Title:       pr.Title,
		Body:        pr.Body,
		State:       pr.State,
		MergeCommit: pr.MergeSHA,
		HeadSHA:     pr.Head.SHA,
		BaseBranch:  pr.Base.Ref,
		HeadBranch:  pr.Head.Ref,
		Merged:      pr.Merged,
		Squashed:    squashed,
		Author:      pr.User.Login,
		MergedAt:    mergedAt,
	}

	return info, nil
}

// GetCommit retrieves information about a commit by SHA.
func (f *Forgejo) GetCommit(ctx context.Context, owner, repo, sha string) (*CommitInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/git/commits/%s", f.baseURL, owner, repo, sha)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if f.token != "" {
		req.Header.Set("Authorization", "token "+f.token)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit %s: %w", sha, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get commit %s: %s (%s)", sha, resp.Status, parseForgejoError(body))
	}

	var commit forgejoCommit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return nil, fmt.Errorf("failed to decode commit response: %w", err)
	}

	parents := make([]string, len(commit.Parents))
	for i, parent := range commit.Parents {
		parents[i] = parent.SHA
	}

	timestamp, _ := time.Parse(time.RFC3339, commit.Commit.Author.Date)

	info := &CommitInfo{
		SHA:       commit.SHA,
		Message:   commit.Commit.Message,
		Author:    commit.Commit.Author.Name,
		Email:     commit.Commit.Author.Email,
		Timestamp: timestamp,
		Parents:   parents,
	}

	return info, nil
}

// ListRecentPRs lists recently merged PRs.
func (f *Forgejo) ListRecentPRs(ctx context.Context, owner, repo string, limit int) ([]*PRInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s/pulls?state=closed&sort=recentupdate&limit=%d", f.baseURL, owner, repo, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if f.token != "" {
		req.Header.Set("Authorization", "token "+f.token)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list PRs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list PRs: %s (%s)", resp.Status, parseForgejoError(body))
	}

	var prs []forgejoPR
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("failed to decode PR list response: %w", err)
	}

	var result []*PRInfo
	for _, pr := range prs {
		if !pr.Merged {
			continue
		}

		mergedAt, _ := time.Parse(time.RFC3339, pr.MergedAt)

		info := &PRInfo{
			Number:      pr.Number,
			Title:       pr.Title,
			State:       pr.State,
			MergeCommit: pr.MergeSHA,
			HeadSHA:     pr.Head.SHA,
			BaseBranch:  pr.Base.Ref,
			HeadBranch:  pr.Head.Ref,
			Merged:      pr.Merged,
			Author:      pr.User.Login,
			MergedAt:    mergedAt,
		}
		result = append(result, info)

		if len(result) >= limit {
			break
		}
	}

	return result, nil
}
