// Package main is the entry point to the server. It reads configuration, sets up logging and error handling,
// handles signals from the OS, and starts and stops the server.
package main

import (
	"canvas/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

func main() {

	os.Exit(start())
}

func start() int {
	logEnv := getStringOrDefaultValue("LOG_ENV", "development")

	log, err := createLogger(logEnv)
	if err != nil {
		fmt.Println("Error setting up the logger", err)
		return 1
	}

	log = log.With(zap.String("release", release))

	defer func() {
		_ = log.Sync()
	}()

	host := getStringOrDefaultValue("HOST", "localhost")
	port := getIntOrDefaultValue("PORT", 8080)

	s := server.New(server.Options{Host: host, Port: port, Log: log})

	// That function makes sure to cancel the returned context if it receives one of the signals we have asked for.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer stop()

	// An error group is something that can run functions in goroutines, wait for all the functions to finish, and return any errors.

	// The Context returned by the errgroup.WithContext function is cancelled the first time any of the functions return an error,
	// which enables us to shut down immediately if there's an error starting up.
	eg, ctx := errgroup.WithContext(ctx)

	// We then call eg.Go with a small function which starts the server and returns any error
	eg.Go(func() error {
		if err := s.Start(); err != nil {
			log.Info("Error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	// After that, we block until the ctx from before is cancelled, by reading from the channel returned by ctx.Done.
	<-ctx.Done()

	// hen it is cancelled (either by a signal or a function passed to eg.Go returning an error),
	// we call Stop on our server in another goroutine passed to the error group.
	eg.Go(func() error {
		if err := s.Stop(); err != nil {
			log.Info("Error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}

}

// Get Env Variables
func getStringOrDefaultValue(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}

func getIntOrDefaultValue(name string, defaultV int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsInt, err := strconv.Atoi(v)
	if err != nil {
		return defaultV
	}
	return vAsInt
}
