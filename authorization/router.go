/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"context"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(ctx context.Context, router *mux.Router, handler API, middleware middleware.APIMiddleware) {

	router.Handle("/signup", handler.signUp(ctx)).Methods(http.MethodPost)

	router.Handle("/u/{confirmationKey:[abcefghijkmnopqrtuvwxyz23479]{16}}",
		handler.confirmUserEmail(ctx)).Methods(http.MethodGet)

	router.Handle("/resetpassword", handler.requestUserPasswordReset(ctx)).Methods(http.MethodPost)

	router.Handle("/oauth/token", handler.createAccessToken(ctx)).Methods(http.MethodPost)
	router.Handle("/oauth/token", handler.refreshAccessToken(ctx)).Methods(http.MethodPut)

	routerSession := router.PathPrefix("/oauth/sessions").Subrouter()
	routerSession.Use(middleware.Authorization)
	routerSession.Handle("", handler.getListActiveSessions(ctx)).Methods(http.MethodGet)
	routerSession.Handle("", handler.revokeRefreshToken(ctx)).Methods(http.MethodDelete)
}
