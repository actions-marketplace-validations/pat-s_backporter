// Package common provides shared CLI flags and utilities.
package common

import (
	"github.com/urfave/cli/v3"

	"codefloe.com/pat-s/backporter/shared/logger"
)

// GlobalFlags are flags available to all commands.
var GlobalFlags = append([]cli.Flag{
	&cli.StringFlag{
		Sources: cli.EnvVars("BACKPORTER_CONFIG"),
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "path to config file",
	},
	&cli.StringFlag{
		Sources: cli.EnvVars("BACKPORTER_REMOTE"),
		Name:    "remote",
		Usage:   "git remote name",
		Value:   "origin",
	},
}, logger.GlobalLoggerFlags...)
