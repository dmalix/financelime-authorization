package rest

import (
	"github.com/dmalix/authorization-service/app/authorization"
	"github.com/dmalix/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.Logger, router *mux.Router, routerV1 *mux.Router, handler authorization.REST, middleware middleware.Middleware) {

	routerV1.Handle("/user/",
		handler.SignUpStep1(logger)).
		Methods(http.MethodPost).
		Headers(headerKeyContentType, headerValueApplicationJson)
	router.Handle("/u/{confirmationKey:[abcefghijkmnopqrtuvwxyz23479]{16}}",
		handler.SignUpStep2(logger)).
		Methods(http.MethodGet)

	routerV1.Handle("/oauth/",
		handler.CreateAccessToken(logger)).
		Methods(http.MethodPost).
		Headers(headerKeyContentType, headerValueApplicationJson)
	routerV1.Handle("/oauth/", handler.RefreshAccessToken(logger)).
		Methods(http.MethodPut).
		Headers(headerKeyContentType, headerValueApplicationJson)

	routerSessions := routerV1.PathPrefix("/sessions").Subrouter()
	routerSessions.Use(middleware.Authorization(logger.Named("middlewareAuthorization")))
	routerSessions.Handle("/",
		handler.GetListActiveSessions(logger)).
		Methods(http.MethodGet)

	routerSession := routerV1.PathPrefix("/session").Subrouter()
	routerSession.Use(middleware.Authorization(logger.Named("middlewareAuthorization")))
	routerSession.Handle("/",
		handler.RevokeRefreshToken(logger)).
		Methods(http.MethodDelete).
		Headers(headerKeyContentType, headerValueApplicationJson)

	routerV1.Handle("/user/",
		handler.ResetUserPasswordStep1(logger)).
		Methods(http.MethodPut).
		Headers(headerKeyContentType, headerValueApplicationJson)
	router.Handle("/p/{confirmationKey:[abcefghijkmnopqrtuvwxyz23479]{16}}",
		handler.ResetUserPasswordStep2(logger)).
		Methods(http.MethodGet)

}
