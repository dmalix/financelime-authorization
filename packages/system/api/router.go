/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"github.com/dmalix/financelime-authorization/packages/authorization"
	"github.com/dmalix/financelime-authorization/packages/system"
	"github.com/dmalix/financelime-authorization/utils/router"
	"net/http"
)

func Router(mux *http.ServeMux, service system.Service, middleware authorization.APIMiddleware) {

	handler := NewHandler(service)

	mux.Handle("/v1/dist",
		router.Group(
			router.EndPoint(router.Point{Method: http.MethodGet, Handler: handler.Dist()}),
			middleware.RequestID,
		))
}
