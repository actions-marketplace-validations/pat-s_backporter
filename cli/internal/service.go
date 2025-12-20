// Package internal provides CLI internal utilities.
package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/cli/internal/config"
	"codefloe.com/pat-s/backporter/pkg/backport"
	pkgconfig "codefloe.com/pat-s/backporter/pkg/config"
	"codefloe.com/pat-s/backporter/pkg/forge"
	"codefloe.com/pat-s/backporter/pkg/git"
)

// CreateService creates a backport service from CLI context.
func CreateService(_ context.Context, c *cli.Command) (*backport.Service, error) {
	cfg, err := config.GetConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Open repository.
	repo, err := git.OpenCurrent()
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get remote URL and parse owner/repo.
	remote := c.String("remote")
	if remote == "" {
		remote = cfg.Remote
	}

	remoteURL, err := repo.RemoteURL(remote)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote URL: %w", err)
	}

	owner, repoName, err := git.ParseRemoteURL(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	log.Debug().Str("owner", owner).Str("repo", repoName).Msg("parsed repository info")

	// Create forge client if configured.
	var f forge.Forge
	if cfg.ForgeType != "" {
		token := getForgeToken(cfg.ForgeType)
		opts := forge.NewOptions{
			ForgejoURL: cfg.ForgejoURL,
		}
		f, err = forge.NewWithOptions(cfg.ForgeType, token, opts)
		if err != nil {
			log.Warn().Err(err).Msg("failed to create forge client")
		} else {
			log.Debug().Str("forge", cfg.ForgeType).Msg("forge client created")
		}
	}

	return backport.NewService(repo, f, cfg, owner, repoName), nil
}

// getForgeToken retrieves the token for the specified forge type from environment.
func getForgeToken(forgeType string) string {
	switch forgeType {
	case "github":
		return os.Getenv("GITHUB_TOKEN")
	case "forgejo":
		return os.Getenv("FORGEJO_TOKEN")
	default:
		return ""
	}
}

// GetRepository opens the current git repository.
func GetRepository() (*git.Repository, error) {
	return git.OpenCurrent()
}

// CreateServiceWithDetails creates a backport service and returns additional details.
// Returns: service, config, forge client, owner, repo name, error.
func CreateServiceWithDetails(_ context.Context, c *cli.Command) (
	*backport.Service,
	*pkgconfig.Config,
	forge.Forge,
	string,
	string,
	error,
) {
	cfg, err := config.GetConfig(c)
	if err != nil {
		return nil, nil, nil, "", "", fmt.Errorf("failed to load config: %w", err)
	}

	// Open repository.
	repo, err := git.OpenCurrent()
	if err != nil {
		return nil, nil, nil, "", "", fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get remote URL and parse owner/repo.
	remote := c.String("remote")
	if remote == "" {
		remote = cfg.Remote
	}

	remoteURL, err := repo.RemoteURL(remote)
	if err != nil {
		return nil, nil, nil, "", "", fmt.Errorf("failed to get remote URL: %w", err)
	}

	owner, repoName, err := git.ParseRemoteURL(remoteURL)
	if err != nil {
		return nil, nil, nil, "", "", fmt.Errorf("failed to parse remote URL: %w", err)
	}

	log.Debug().Str("owner", owner).Str("repo", repoName).Msg("parsed repository info")

	// Create forge client if configured.
	var f forge.Forge
	if cfg.ForgeType != "" {
		token := getForgeToken(cfg.ForgeType)
		opts := forge.NewOptions{
			ForgejoURL: cfg.ForgejoURL,
		}
		f, err = forge.NewWithOptions(cfg.ForgeType, token, opts)
		if err != nil {
			return nil, nil, nil, "", "", fmt.Errorf("failed to create forge client: %w", err)
		}
		log.Debug().Str("forge", cfg.ForgeType).Msg("forge client created")
	} else {
		return nil, nil, nil, "", "", fmt.Errorf("forge_type must be configured for CI mode")
	}

	svc := backport.NewService(repo, f, cfg, owner, repoName)
	return svc, cfg, f, owner, repoName, nil
}
