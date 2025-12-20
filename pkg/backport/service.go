package backport

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"codefloe.com/pat-s/backporter/pkg/config"
	"codefloe.com/pat-s/backporter/pkg/forge"
	"codefloe.com/pat-s/backporter/pkg/git"
	"codefloe.com/pat-s/backporter/shared/version"
)

// Service orchestrates backport operations.
type Service struct {
	repo   *git.Repository
	forge  forge.Forge
	config *config.Config
	cache  *Cache
	owner  string
	repoN  string
}

// NewService creates a new backport service.
func NewService(repo *git.Repository, f forge.Forge, cfg *config.Config, owner, repoName string) *Service {
	cachePath := cfg.Cache.Path
	if !cfg.Cache.Enabled {
		cachePath = ""
	}

	return &Service{
		repo:   repo,
		forge:  f,
		config: cfg,
		cache:  NewCache(cachePath),
		owner:  owner,
		repoN:  repoName,
	}
}

// BackportOptions contains options for backport operations.
type BackportOptions struct {
	TargetBranch string
	DryRun       bool
}

// BackportResult contains the result of a backport operation.
type BackportResult struct {
	OriginalSHA  string
	BackportSHA  string
	TargetBranch string
	PRNumber     int
	Success      bool
	HasConflict  bool
	Message      string
}

// BackportCommit backports a single commit to the target branch.
func (s *Service) BackportCommit(_ context.Context, sha string, opts BackportOptions) (*BackportResult, error) {
	log.Debug().Str("sha", sha).Str("target", opts.TargetBranch).Msg("backporting commit")

	// Verify the commit exists.
	fullSHA, err := s.repo.GetCommitSHA(sha)
	if err != nil {
		return nil, fmt.Errorf("commit not found: %w", err)
	}

	// Check for uncommitted changes.
	hasChanges, err := s.repo.HasUncommittedChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to check for uncommitted changes: %w", err)
	}
	if hasChanges {
		return nil, fmt.Errorf("repository has uncommitted changes, please commit or stash them first")
	}

	// Store original branch.
	originalBranch, err := s.repo.CurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	// Verify target branch exists.
	exists, err := s.repo.BranchExists(opts.TargetBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to check target branch: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("target branch %s does not exist", opts.TargetBranch)
	}

	if opts.DryRun {
		log.Info().Msg("dry-run mode, not making changes")
		return &BackportResult{
			OriginalSHA:  fullSHA,
			TargetBranch: opts.TargetBranch,
			Success:      true,
			Message:      "dry-run: would backport commit",
		}, nil
	}

	// Checkout target branch.
	log.Debug().Str("branch", opts.TargetBranch).Msg("checking out target branch")
	if err := git.CheckoutBranch(opts.TargetBranch); err != nil {
		return nil, err
	}

	// Track whether we should return to original branch.
	shouldCheckoutBack := true

	// Ensure we return to original branch on error (unless conflict).
	defer func() {
		if shouldCheckoutBack && originalBranch != "" {
			_ = git.CheckoutBranch(originalBranch)
		}
	}()

	// Perform cherry-pick.
	log.Debug().Str("sha", fullSHA).Msg("cherry-picking commit")
	result, err := git.CherryPick(fullSHA)
	if err != nil {
		return nil, err
	}

	if result.HasConflict {
		// Don't switch back to original branch - user needs to resolve conflicts.
		shouldCheckoutBack = false
		return &BackportResult{
			OriginalSHA:  fullSHA,
			TargetBranch: opts.TargetBranch,
			Success:      false,
			HasConflict:  true,
			Message:      result.Message,
		}, nil
	}

	// Get the new commit SHA.
	newSHA, err := git.GetCurrentCommitSHA()
	if err != nil {
		return nil, fmt.Errorf("failed to get new commit SHA: %w", err)
	}

	// Amend commit message with backport signature.
	originalMessage, err := s.repo.GetCommitMessage(newSHA)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit message: %w", err)
	}

	signature := version.SignatureMessage(fullSHA)
	newMessage := fmt.Sprintf("%s\n\n%s", originalMessage, signature)

	if err := git.AmendCommitMessage(newMessage); err != nil {
		return nil, fmt.Errorf("failed to amend commit message: %w", err)
	}

	// Get final SHA after amend.
	finalSHA, err := git.GetCurrentCommitSHA()
	if err != nil {
		return nil, fmt.Errorf("failed to get final commit SHA: %w", err)
	}

	// Cache the result.
	if s.cache != nil && s.config.Cache.Enabled {
		entry := CacheEntry{
			OriginalSHA:  fullSHA,
			BackportSHA:  finalSHA,
			TargetBranch: opts.TargetBranch,
			Timestamp:    time.Now(),
			Message:      originalMessage,
		}
		if err := s.cache.Add(entry); err != nil {
			log.Warn().Err(err).Msg("failed to cache backport entry")
		}
	}

	log.Debug().Str("sha", finalSHA).Msg("commit successfully backported")

	return &BackportResult{
		OriginalSHA:  fullSHA,
		BackportSHA:  finalSHA,
		TargetBranch: opts.TargetBranch,
		Success:      true,
		Message:      "commit successfully backported",
	}, nil
}

// BackportPR backports a PR's merge commit to the target branch.
func (s *Service) BackportPR(ctx context.Context, prNumber int, opts BackportOptions) (*BackportResult, error) {
	if s.forge == nil {
		return nil, fmt.Errorf("forge not configured, cannot backport PR")
	}

	log.Debug().Int("pr", prNumber).Str("target", opts.TargetBranch).Msg("backporting PR")

	// Fetch PR information.
	prInfo, err := s.forge.GetPR(ctx, s.owner, s.repoN, prNumber)
	if err != nil {
		return nil, err
	}

	// Check if PR was squash merged.
	if !prInfo.IsSquashMerge() {
		return nil, fmt.Errorf("PR #%d was not squash merged - please backport individual commits instead", prNumber)
	}

	// Backport the merge commit.
	result, err := s.BackportCommit(ctx, prInfo.MergeCommit, opts)
	if err != nil {
		return nil, err
	}

	result.PRNumber = prNumber

	// Update cache with PR number.
	if s.cache != nil && s.config.Cache.Enabled && result.Success {
		entries := s.cache.FindByOriginalSHA(result.OriginalSHA)
		if len(entries) > 0 {
			// Update the last entry with PR number.
			lastIdx := len(s.cache.entries) - 1
			s.cache.entries[lastIdx].PRNumber = prNumber
			_ = s.cache.save()
		}
	}

	return result, nil
}

// ListBackports returns the list of cached backport operations.
func (s *Service) ListBackports() []CacheEntry {
	if s.cache == nil {
		return nil
	}
	return s.cache.List()
}

// ClearCache clears the backport cache.
func (s *Service) ClearCache() error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Clear()
}
