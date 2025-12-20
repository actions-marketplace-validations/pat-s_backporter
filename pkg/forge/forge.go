package forge

import (
	"context"
	"fmt"
	"os"
)

// Forge is the interface for interacting with git forges.
type Forge interface {
	// GetPR retrieves information about a pull request by number.
	GetPR(ctx context.Context, owner, repo string, number int) (*PRInfo, error)

	// GetCommit retrieves information about a commit by SHA.
	GetCommit(ctx context.Context, owner, repo, sha string) (*CommitInfo, error)

	// ListRecentPRs lists recently merged PRs.
	ListRecentPRs(ctx context.Context, owner, repo string, limit int) ([]*PRInfo, error)

	// Name returns the name of the forge.
	Name() string
}

// NewOptions holds options for creating a forge client.
type NewOptions struct {
	ForgejoURL string // Required for Forgejo forge type
}

// New creates a new forge client based on the forge type.
func New(forgeType, token string) (Forge, error) {
	return NewWithOptions(forgeType, token, NewOptions{})
}

// NewWithOptions creates a new forge client with additional options.
func NewWithOptions(forgeType, token string, opts NewOptions) (Forge, error) {
	switch forgeType {
	case "github":
		return NewGitHub(token), nil
	case "forgejo":
		// Forgejo requires a base URL - check options first, then environment.
		baseURL := opts.ForgejoURL
		if baseURL == "" {
			baseURL = os.Getenv("FORGEJO_URL")
		}
		if baseURL == "" {
			return nil, fmt.Errorf("FORGEJO_URL not configured (set in config file or FORGEJO_URL environment variable)")
		}
		return NewForgejo(baseURL, token), nil
	default:
		return nil, fmt.Errorf("unknown forge type: %s", forgeType)
	}
}
