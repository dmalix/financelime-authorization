/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"github.com/dmalix/financelime-rest-api/packages/authorization/domain"
	"github.com/dmalix/financelime-rest-api/utils/middleware"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, service domain.AccountService) {

	handler := NewHandler(service)

	mux.Handle("/authorization/signup", middleware.RequestID(handler.SignUp()))

}
