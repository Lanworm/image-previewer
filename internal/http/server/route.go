package server

import (
	"github.com/Lanworm/image-previewe/internal/http/server/httphandler"
)

const baseContentType = "application/json"

func (s *Server) RegisterRoutes(handler *httphandler.Handler) {
	s.AddRoute("/fill/", ContentType(baseContentType, Method("GET", handler.ResizeHandler)))
}
