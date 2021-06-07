/* Copyright © 2021. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package rest

import (
	"github.com/dmalix/financelime-authorization/app/information"
	"github.com/dmalix/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.Logger, router *mux.Router, handler information.REST, _ middleware.Middleware) {
	router.Handle("/version", handler.Version(logger)).Methods(http.MethodGet)
}
