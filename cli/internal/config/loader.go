// Package config provides CLI-specific configuration loading.
package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/pkg/config"
)

// Load loads configuration from global and repo-local config files.
func Load(c *cli.Command) (*config.Config, error) {
	cfg := config.DefaultConfig()

	// Load global config first.
	globalPath := config.GlobalConfigPath()
	if globalPath != "" {
		if _, err := os.Stat(globalPath); err == nil {
			globalCfg, err := config.LoadFromFile(globalPath)
			if err != nil {
				log.Debug().Err(err).Str("path", globalPath).Msg("failed to load global config")
			} else {
				log.Debug().Str("path", globalPath).Msg("loaded global config")
				cfg.Merge(globalCfg)
			}
		}
	}

	// Load repo-local config (overrides global).
	repoPath := config.RepoConfigPath()
	if _, err := os.Stat(repoPath); err == nil {
		repoCfg, err := config.LoadFromFile(repoPath)
		if err != nil {
			log.Debug().Err(err).Str("path", repoPath).Msg("failed to load repo config")
		} else {
			log.Debug().Str("path", repoPath).Msg("loaded repo config")
			cfg.Merge(repoCfg)
		}
	}

	// Override with explicit config file if provided.
	if configPath := c.String("config"); configPath != "" {
		explicitCfg, err := config.LoadFromFile(configPath)
		if err != nil {
			return nil, err
		}
		log.Debug().Str("path", configPath).Msg("loaded explicit config")
		cfg.Merge(explicitCfg)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Warn if forge type is not set.
	if cfg.ForgeType == "" {
		log.Warn().Msg("forge_type not configured - PR features will be unavailable")
	}

	return cfg, nil
}

// ApplyToFlags applies config values to CLI flags if they haven't been explicitly set.
func ApplyToFlags(c *cli.Command, cfg *config.Config) error {
	// Only apply if the flag hasn't been explicitly set.
	if !c.IsSet("remote") && cfg.Remote != "" {
		if err := c.Set("remote", cfg.Remote); err != nil {
			return err
		}
	}

	return nil
}

// GetConfig retrieves the current configuration from context.
func GetConfig(c *cli.Command) (*config.Config, error) {
	return Load(c)
}
