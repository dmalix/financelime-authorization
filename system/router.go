/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package system

import (
	"context"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(ctx context.Context, router *mux.Router, handler API, _ middleware.APIMiddleware) {
	router.Handle("/version", handler.version(ctx)).Methods(http.MethodGet)
}
