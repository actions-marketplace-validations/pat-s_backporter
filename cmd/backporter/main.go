// Package main is the entry point for the backporter CLI.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle signals for graceful shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info().Msg("termination signal received, shutting down")
		cancel()
	}()

	app := newApp()
	if err := app.Run(ctx, os.Args); err != nil {
		cancel()
		log.Fatal().Err(err).Msg("error running backporter")
	}
	cancel()
}
