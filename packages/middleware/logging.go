package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func (mw *middleware) Logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				panicStatus := recover()
				if panicStatus != nil {
					logger.DPanic("Recover has panic status",
						zap.Any("error", panicStatus),
						zap.String("stack", string(debug.Stack())))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)

			remoteAddr, err := getRemoteAddr(r.Context())
			if err != nil {
				logger.DPanic("failed to get remoteAddr", zap.Error(err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}

			requestID, requestIDKey, err := getRequestID(r.Context())
			if err != nil {
				logger.DPanic("failed to get requestID", zap.Error(err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}

			inner.ServeHTTP(wrapped, r)

			if wrapped.status >= 100 && wrapped.status < 200 {
				logger.Warn("Informational",
					zap.String(ContextKeyRemoteAddr, remoteAddr),
					zap.String(requestIDKey, requestID),
					zap.Int("status", wrapped.status),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
			} else if wrapped.status >= 200 && wrapped.status < 300 {
				logger.Info("Success",
					zap.String(ContextKeyRemoteAddr, remoteAddr),
					zap.String(requestIDKey, requestID),
					zap.Int("status", wrapped.status),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
			} else if wrapped.status >= 300 && wrapped.status < 400 {
				logger.Warn("Redirection",
					zap.String(ContextKeyRemoteAddr, remoteAddr),
					zap.String(requestIDKey, requestID),
					zap.Int("status", wrapped.status),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
			} else if wrapped.status >= 400 && wrapped.status < 500 {
				logger.Error("Client Error",
					zap.String(ContextKeyRemoteAddr, remoteAddr),
					zap.String(requestIDKey, requestID),
					zap.Int("status", wrapped.status),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
			} else {
				logger.Error("Server Error",
					zap.String(ContextKeyRemoteAddr, remoteAddr),
					zap.String(requestIDKey, requestID),
					zap.Int("status", wrapped.status),
					zap.String("method", r.Method),
					zap.String("path", r.URL.EscapedPath()),
					zap.Duration("duration", time.Since(start)))
			}
		})
	}
}
