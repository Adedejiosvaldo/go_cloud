package server

import "canvas/handlers"

func (s *Server) SetupRoutes() {
	handlers.Health(s.mux)
}
