// Package git provides git operations using go-git.
package git

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repository wraps go-git repository operations.
type Repository struct {
	repo *gogit.Repository
}

// Open opens an existing git repository.
func Open(path string) (*Repository, error) {
	repo, err := gogit.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &Repository{repo: repo}, nil
}

// OpenCurrent opens the git repository in the current directory or any parent.
func OpenCurrent() (*Repository, error) {
	repo, err := gogit.PlainOpenWithOptions(".", &gogit.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &Repository{repo: repo}, nil
}

// RemoteURL returns the URL of the specified remote.
func (r *Repository) RemoteURL(name string) (string, error) {
	remote, err := r.repo.Remote(name)
	if err != nil {
		return "", fmt.Errorf("failed to get remote %s: %w", name, err)
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("remote %s has no URLs", name)
	}

	return urls[0], nil
}

// ParseRemoteURL parses a git remote URL and extracts owner and repo.
func ParseRemoteURL(url string) (owner, repo string, err error) {
	// Handle SSH URLs: git@github.com:owner/repo.git
	if strings.HasPrefix(url, "git@") {
		parts := strings.Split(url, ":")
		if len(parts) != 2 { //nolint:mnd
			return "", "", fmt.Errorf("invalid SSH URL format: %s", url)
		}
		path := strings.TrimSuffix(parts[1], ".git")
		pathParts := strings.Split(path, "/")
		if len(pathParts) != 2 { //nolint:mnd
			return "", "", fmt.Errorf("invalid SSH URL path: %s", url)
		}
		return pathParts[0], pathParts[1], nil
	}

	// Handle HTTPS URLs: https://github.com/owner/repo.git
	re := regexp.MustCompile(`https?://[^/]+/([^/]+)/([^/]+?)(?:\.git)?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 { //nolint:mnd
		return "", "", fmt.Errorf("invalid HTTPS URL format: %s", url)
	}

	return matches[1], matches[2], nil
}

// CurrentBranch returns the name of the current branch.
func (r *Repository) CurrentBranch() (string, error) {
	head, err := r.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if !head.Name().IsBranch() {
		return "", fmt.Errorf("HEAD is not pointing to a branch")
	}

	return head.Name().Short(), nil
}

// HasUncommittedChanges checks if there are uncommitted changes.
// This only checks for modified or staged files, not untracked files.
func (r *Repository) HasUncommittedChanges() (bool, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get status: %w", err)
	}

	// Check for modified or staged files only (ignore untracked files).
	for _, fileStatus := range status {
		// Skip untracked files (Worktree: Untracked, Staging: Untracked).
		if fileStatus.Worktree == gogit.Untracked && fileStatus.Staging == gogit.Untracked {
			continue
		}
		// Any other status means there are uncommitted changes.
		return true, nil
	}

	return false, nil
}

// BranchExists checks if a branch exists.
func (r *Repository) BranchExists(name string) (bool, error) {
	_, err := r.repo.Reference(plumbing.NewBranchReferenceName(name), true)
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListBranches returns a list of branch names.
func (r *Repository) ListBranches() ([]string, error) {
	iter, err := r.repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	var branches []string
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})

	return branches, err
}

// GetCommitSHA returns the SHA of a commit reference (branch name, tag, or SHA).
func (r *Repository) GetCommitSHA(ref string) (string, error) {
	hash, err := r.repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return "", fmt.Errorf("failed to resolve %s: %w", ref, err)
	}

	return hash.String(), nil
}

// GetCommitMessage returns the commit message for a given SHA.
func (r *Repository) GetCommitMessage(sha string) (string, error) {
	hash := plumbing.NewHash(sha)
	commit, err := r.repo.CommitObject(hash)
	if err != nil {
		return "", fmt.Errorf("failed to get commit %s: %w", sha, err)
	}

	return commit.Message, nil
}

// Inner returns the underlying go-git repository.
func (r *Repository) Inner() *gogit.Repository {
	return r.repo
}
