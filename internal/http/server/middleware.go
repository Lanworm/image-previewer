package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lanworm/OTUS_GO/final_project/internal/logger"
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

func ContentType(t string, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", t)
		next(writer, request)
	}
}

func Method(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte(fmt.Sprintf(
				"handler %s support method %s",
				request.URL.Path,
				method,
			)))

			return
		}

		next(writer, request)
	}
}
