package restv1

import (
	"github.com/dmalix/authorization-service/app/authorization"
	"github.com/dmalix/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func RouterV1(logger *zap.Logger, router *mux.Router, handler authorization.RestV1, middleware middleware.Middleware) {

	router.Handle("/user/",
		handler.SignUpStep1(logger)).
		Methods(http.MethodPost).
		Headers(headerKeyContentType, headerValueApplicationJson)

	router.Handle("/oauth/",
		handler.CreateAccessToken(logger)).
		Methods(http.MethodPost).
		Headers(headerKeyContentType, headerValueApplicationJson)
	router.Handle("/oauth/", handler.RefreshAccessToken(logger)).
		Methods(http.MethodPut).
		Headers(headerKeyContentType, headerValueApplicationJson)

	routerSessions := router.PathPrefix("/sessions").Subrouter()
	routerSessions.Use(middleware.Authorization(logger.Named("middlewareAuthorization")))
	routerSessions.Handle("/",
		handler.GetListActiveSessions(logger)).
		Methods(http.MethodGet)

	routerSession := router.PathPrefix("/session").Subrouter()
	routerSession.Use(middleware.Authorization(logger.Named("middlewareAuthorization")))
	routerSession.Handle("/",
		handler.RevokeRefreshToken(logger)).
		Methods(http.MethodDelete).
		Headers(headerKeyContentType, headerValueApplicationJson)

	router.Handle("/user/",
		handler.ResetUserPasswordStep1(logger)).
		Methods(http.MethodPut).
		Headers(headerKeyContentType, headerValueApplicationJson)

	router.Handle("/u/{confirmationKey:[abcefghijkmnopqrtuvwxyz23479]{16}}",
		handler.SignUpStep2(logger)).
		Methods(http.MethodGet)

	router.Handle("/p/{confirmationKey:[abcefghijkmnopqrtuvwxyz23479]{16}}",
		handler.ResetUserPasswordStep2(logger)).
		Methods(http.MethodGet)

}
