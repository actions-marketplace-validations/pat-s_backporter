package common

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/cli/internal/config"
	"codefloe.com/pat-s/backporter/cli/setup"
	"codefloe.com/pat-s/backporter/shared/logger"
)

// Before is the global before hook that sets up logging and loads config.
func Before(ctx context.Context, c *cli.Command) (context.Context, error) {
	if err := logger.SetupGlobalLogger(ctx, c); err != nil {
		return ctx, err
	}

	log.Debug().Str("version", c.Root().Version).Msg("backporter starting")

	// Check if we should prompt for config creation.
	if setup.ShouldPromptForConfig() && !logger.IsCI() && c.String("config") == "" {
		if err := setup.PromptForConfigCreation(); err != nil {
			log.Warn().Err(err).Msg("failed to create config")
		}
	}

	// Load configuration.
	cfg, err := config.Load(c)
	if err != nil {
		log.Warn().Err(err).Msg("failed to load config, using defaults")
	} else if cfg != nil {
		// Apply config values to CLI flags if not already set.
		if err := config.ApplyToFlags(c, cfg); err != nil {
			log.Warn().Err(err).Msg("failed to apply config to flags")
		}
	}

	return ctx, nil
}
