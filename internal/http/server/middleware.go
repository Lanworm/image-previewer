package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lanworm/image-previewer/internal/logger"
)

type Logging struct {
	logger *logger.Logger
}

func NewLogMiddleware(logger *logger.Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (l *Logging) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.logger.ServerLog(fmt.Sprintf(
			"%s [%s] %s %s %s \"%s\"",
			r.RemoteAddr,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			r.Proto,
			r.Header.Get("User-Agent"),
		))
		next.ServeHTTP(w, r)
	})
}

type Recovery struct {
	logger *logger.Logger
}

func NewRecoveryMiddleware(logger *logger.Logger) *Recovery {
	return &Recovery{
		logger: logger,
	}
}

func (rc *Recovery) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if res := recover(); res != nil {
				rc.logger.ServerLog(fmt.Sprintf("[RECOVERY] %s", res))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
