/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package router

import (
	"net/http"
)

func exampleRouter(mux *http.ServeMux) {
	mux.Handle("/test",
		Group(
			EndPoint(
				Point{Method: http.MethodPost, Handler: handlerMock1()},
				Point{Method: http.MethodPut, Handler: handlerMock2()}),
		))
}
