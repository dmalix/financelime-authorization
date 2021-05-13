package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status       int
	headerStatus bool
}

func replacementResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
	}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.headerStatus {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.headerStatus = true

	return
}

func (middleware *Middleware) Logging() func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				panicStatus := recover()
				if panicStatus != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Println(
						"err", panicStatus,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			replacement := replacementResponseWriter(w)

			inner.ServeHTTP(replacement, r)

			log.Println(
				"status", replacement.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		})
	}
}
