// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"canvas/storage"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	address  string
	database *storage.Database
	mux      chi.Router
	server   *http.Server
	log      *zap.Logger
}

type Options struct {
	Database *storage.Database
	Host     string
	Port     int
	Log      *zap.Logger
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		address:  address,
		mux:      mux,
		database: opts.Database,
		log:      opts.Log,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       5 * time.Second,
			WriteTimeout:      5 * time.Second,
		},
	}
}

func (s *Server) Start() error {

	if err := s.database.Connect(); err != nil {
		return fmt.Errorf("error connecting to databse: %w", err)
	}
	s.SetupRoutes()

	// fmt.Println("Starting on", s.address)
	// using zapper
	s.log.Info("Starting the server", zap.String("address", s.address))

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("Error Starting Server: %w", err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.log.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Error stopping server: %w", err)
	}

	return nil
}
