package middleware

import (
	"context"
	"go.uber.org/zap"
	"html"
	"net/http"
)

func (mw *middleware) RequestID(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			const isRequired = true
			const toCheck = true
			var requestID string

			remoteAddr, remoteAddrKey, err := getRemoteAddr(r.Context())
			if err != nil {
				logger.DPanic("failed to get remoteAddr", zap.Error(err))
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}

			if isRequired {
				requestID = html.EscapeString(r.Header.Get("request-id"))
				if requestID == "" {
					if r.Method == http.MethodGet { // Try get RequestID value from Url Param 'rid' (for confirm user email)
						rid, ok := r.URL.Query()["rid"]
						if !ok || len(rid[0]) < 1 {
							logger.Error("the header 'request-id' and the Url-param 'rid' not found but required")
							return
						}
						requestID = rid[0]
					} else {
						logger.Error("the 'request-id' header not found but required", zap.String(remoteAddrKey, remoteAddr))
						http.Error(w, statusMessageBadRequest, http.StatusBadRequest)
						return
					}
				}
				if toCheck {
					if len(requestID) != 29 {
						logger.Error("the 'request-id' value is invalid",
							zap.String("requestID", requestID), zap.String(remoteAddrKey, remoteAddr))
						http.Error(w, statusMessageBadRequest, http.StatusBadRequest)
						return
					}
				}
			}

			ctx := context.WithValue(r.Context(), ContextKeyRequestID, requestID)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
