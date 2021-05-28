package middleware

import (
	"context"
	"github.com/dmalix/financelime-authorization/packages/jwt"
	"go.uber.org/zap"
	"html"
	"net/http"
	"strings"
)

func (mw *middleware) Authorization(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var (
				err           error
				authorization string
				jwtTokenArr   []string
				jwtData       jwt.JsonWebToken
			)

			remoteAddr, remoteAddrKey, err := getRemoteAddr(r.Context())
			if err != nil {
				logger.DPanic("failed to get remoteAddr", zap.Error(err))
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}

			requestID, requestIDKey, err := getRequestID(r.Context())
			if err != nil {
				logger.DPanic("failed to get requestID", zap.Error(err))
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}

			// Get an authorization token from the header

			authorization = html.EscapeString(r.Header.Get("authorization"))
			if authorization == "" {
				logger.Error("the 'authorization' header not found",
					zap.String(requestIDKey, requestID), zap.String(remoteAddrKey, remoteAddr))
				http.Error(w, statusMessageUnauthorized, http.StatusUnauthorized)
				return
			}

			// Validate Token and data extract for identification

			jwtTokenArr = strings.Split(strings.TrimSpace(html.EscapeString(authorization)), " ")
			if len(jwtTokenArr) != 2 {
				logger.Error("the 'authorization' header is invalid",
					zap.String("error", "want 'bearer token'"),
					zap.String(ContextKeyRequestID, requestID), zap.String(remoteAddrKey, remoteAddr))
				http.Error(w, statusMessageUnauthorized, http.StatusUnauthorized)
				return
			}

			if strings.ToLower(jwtTokenArr[0]) != "bearer" {
				logger.Error("the 'authorization' header is invalid",
					zap.String("error", "got 'xxx token' want 'bearer token'"),
					zap.String(ContextKeyRequestID, requestID), zap.String(remoteAddrKey, remoteAddr))
				http.Error(w, statusMessageUnauthorized, http.StatusUnauthorized)
				return
			}

			jwtData, err = mw.jwt.VerifyToken(jwtTokenArr[1])
			if err != nil {
				logger.Error("the jwt-token is not valid",
					zap.String("jwt", jwtTokenArr[1]), zap.Error(err),
					zap.String(ContextKeyRequestID, requestID), zap.String(remoteAddrKey, remoteAddr))
				http.Error(w, statusMessageUnauthorized, http.StatusUnauthorized)
				return
			}

			if jwtData.Payload.Purpose != "access" {
				logger.Error("the jwt-token is not valid",
					zap.String("jwt", jwtTokenArr[1]),
					zap.String("error", "'Payload.Purpose' must be 'access'"),
					zap.String(ContextKeyRequestID, requestID), zap.String(remoteAddrKey, remoteAddr))
				http.Error(w, statusMessageUnauthorized+" [%s]", http.StatusUnauthorized)
				return
			}

			logger.Info("Authorization was successful",
				zap.String(remoteAddrKey, remoteAddr), zap.String(ContextKeyRequestID, requestID),
				zap.String(ContextKeyPublicSessionID, jwtData.Payload.PublicSessionID))

			ctx := context.WithValue(r.Context(), ContextKeyPublicSessionID, jwtData.Payload.PublicSessionID)
			ctx = context.WithValue(ctx, ContextKeyEncryptedJWTData, jwtData.Payload.Data)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
