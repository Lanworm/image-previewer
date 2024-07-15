package server

import (
	"github.com/Lanworm/image-previewer/internal/http/server/httphandler"
)

func (s *Server) RegisterRoutes(handler *httphandler.Handler) {
	s.AddRoute("/fill/{width}/{height}/{url:.*}", handler.ResizeHandler)
}
