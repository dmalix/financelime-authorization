/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package router

import (
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/responder"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"log"
	"net/http"
)

type Point struct {
	Method  string
	Handler http.Handler
}

func EndPoint(points ...Point) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			handler http.Handler
			result  bool
		)

		for _, point := range points {
			if r.Method == point.Method {
				result = true
				handler = point.Handler
				break
			}
		}

		if result != true {
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Methods Not Allowed [%s]",
					responder.Message(r)))
			http.Error(w, "405 Methods Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if handler == nil {
			log.Printf("FATAL [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Handler Not Found [%s]",
					responder.Message(r)))
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
