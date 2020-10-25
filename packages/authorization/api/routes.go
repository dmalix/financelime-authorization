/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"fmt"
	"github.com/dmalix/financelime-rest-api/packages/authorization/domain"
	"github.com/dmalix/financelime-rest-api/utils/responder"
	"log"
	"net/http"
)

type point struct {
	method  string
	handler http.Handler
}

func Router(mux *http.ServeMux, service domain.AccountService, middleware ...func(http.Handler) http.Handler) {

	handler := NewHandler(service)

	mux.Handle("/authorization/signup",
		middlewareGroup(
			end(point{http.MethodPost, handler.SignUp()}),
			middleware,
		))
}

func middlewareGroup(h http.Handler, middleware []func(http.Handler) http.Handler) http.Handler {

	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func end(points ...point) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var handler http.Handler
		var result bool

		for _, point := range points {
			if r.Method == point.method {
				result = true
				handler = point.handler
				break
			}
		}

		if result != true {
			label := "226Pi8rl"
			log.Printf("ERROR [%s: %s]", label,
				fmt.Sprintf("Methods Not Allowed [%s]",
					responder.Message(r)))
			http.Error(w, "405 Methods Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if handler == nil {
			label := "bq31WdVJ"
			log.Printf("FATAL [%s: %s]", label,
				fmt.Sprintf("Handler Not Found [%s]",
					responder.Message(r)))
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
