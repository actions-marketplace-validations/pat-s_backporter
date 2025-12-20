// Package logger provides logging setup for the application.
package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// GlobalLoggerFlags returns the global logger flags.
var GlobalLoggerFlags = []cli.Flag{
	&cli.StringFlag{
		Sources: cli.EnvVars("BACKPORTER_LOG_LEVEL"),
		Name:    "log-level",
		Usage:   "set logging level",
		Value:   "info",
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("BACKPORTER_PRETTY"),
		Name:    "pretty",
		Usage:   "enable pretty-printed debug output",
		Value:   isInteractiveTerminal(),
	},
	&cli.BoolFlag{
		Sources: cli.EnvVars("BACKPORTER_NOCOLOR"),
		Name:    "nocolor",
		Usage:   "disable colored debug output",
		Value:   !isInteractiveTerminal(),
	},
}

// SetupGlobalLogger configures the global logger based on CLI flags.
func SetupGlobalLogger(_ context.Context, c *cli.Command) error {
	logLevel := c.String("log-level")
	pretty := c.Bool("pretty")
	noColor := c.Bool("nocolor")

	var out io.Writer = os.Stderr

	log.Logger = zerolog.New(out).With().Timestamp().Logger()

	if pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     out,
				NoColor: noColor,
			},
		)
	}

	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("unknown logging level: %s", logLevel)
	}
	zerolog.SetGlobalLevel(lvl)

	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		log.Logger = log.With().Caller().Logger()
	}

	return nil
}

// isInteractiveTerminal returns true if stdout is an interactive terminal.
func isInteractiveTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// IsCI returns true if running in CI environment.
func IsCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}
