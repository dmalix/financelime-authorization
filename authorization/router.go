/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"github.com/dmalix/financelime-authorization/utils/router"
	"net/http"
)

func Router(mux *http.ServeMux, handler API, middleware APIMiddleware) {

	mux.Handle("/v1/signup",
		router.Group(
			router.EndPoint(router.Point{Method: http.MethodPost, Handler: handler.signUp()})))

	mux.Handle("/u/",
		router.Group(
			router.EndPoint(router.Point{Method: http.MethodGet, Handler: handler.confirmUserEmail()}),
		))

	mux.Handle("/v1/resetpassword",
		router.Group(
			router.EndPoint(router.Point{Method: http.MethodPost, Handler: handler.requestUserPasswordReset()}),
		))

	mux.Handle("/v1/oauth/token",
		router.Group(
			router.EndPoint(
				router.Point{Method: http.MethodPost, Handler: handler.createAccessToken()},
				router.Point{Method: http.MethodPut, Handler: handler.refreshAccessToken()}),
		))

	mux.Handle("/v1/oauth/sessions",
		router.Group(
			router.EndPoint(
				router.Point{Method: http.MethodGet, Handler: handler.getListActiveSessions()},
				router.Point{Method: http.MethodDelete, Handler: handler.revokeRefreshToken()}),
			middleware.authorization,
		))
}
