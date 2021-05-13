/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package rest

import (
	"github.com/dmalix/financelime-authorization/app/authorization"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.Logger, router *mux.Router, routerV1 *mux.Router, handler authorization.REST, middleware middleware.Middleware) {

	routerV1.Handle("/user/signup", handler.SignUp(logger)).
		Methods(http.MethodPost).Headers(headerKeyContentType, headerValueApplicationJson)

	router.Handle("/u/{confirmationKey:[abcefghijkmnopqrtuvwxyz23479]{16}}",
		handler.ConfirmUserEmail(logger)).Methods(http.MethodGet)

	routerV1.Handle("/user/password", handler.RequestUserPasswordReset(logger)).Methods(http.MethodPost).
		Methods(http.MethodPost).Headers(headerKeyContentType, headerValueApplicationJson)

	routerV1.Handle("/oauth/create", handler.CreateAccessToken(logger)).Methods(http.MethodPost).
		Headers(headerKeyContentType, headerValueApplicationJson)
	routerV1.Handle("/oauth/refresh", handler.RefreshAccessToken(logger)).Methods(http.MethodPut).
		Headers(headerKeyContentType, headerValueApplicationJson)

	routerSession := routerV1.PathPrefix("/session").Subrouter()
	routerSession.Use(middleware.Authorization(logger.Named("middlewareAuthorization")))
	routerSession.Handle("/list", handler.GetListActiveSessions(logger)).Methods(http.MethodGet)
	routerSession.Handle("/remove", handler.RevokeRefreshToken(logger)).Methods(http.MethodDelete).
		Headers(headerKeyContentType, headerValueApplicationJson)
}
