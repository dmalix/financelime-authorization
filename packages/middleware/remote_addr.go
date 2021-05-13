package middleware

import (
	"context"
	"go.uber.org/zap"
	"html"
	"net"
	"net/http"
)

func (mw *middleware) RemoteAddr(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			headerXRealIp := true

			remoteAddr := html.EscapeString(r.Header.Get("X-Real-IP"))
			if remoteAddr == "" {
				headerXRealIp = false
				remoteAddr = r.RemoteAddr
			}

			checkAddrSource := net.ParseIP(remoteAddr)
			if checkAddrSource == nil {
				logger.Error("the remote addr is not valid",
					zap.Bool("XRealIp", headerXRealIp), zap.String(ContextKeyRemoteAddr, remoteAddr))
				http.Error(w, statusMessageBadRequest, http.StatusBadRequest)
				return
			}
			resultRemoteAddr := checkAddrSource.String()

			ctx := context.WithValue(r.Context(), ContextKeyRemoteAddr, resultRemoteAddr)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
