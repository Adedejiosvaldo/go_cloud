// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
}

type Options struct {
	Host string
	Port int
}

func New(opts Options) *Server {
	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		address: address,
		mux:     mux,
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
