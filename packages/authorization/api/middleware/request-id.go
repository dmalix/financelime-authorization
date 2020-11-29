/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package middleware

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func (m *Middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const (
			isRequired = true
			toCheck    = true
		)

		if isRequired {
			requestID := r.Header.Get("request-id")
			if len(requestID) == 0 {
				log.Printf("ERROR [%s: %s]", "UV9WNZfq",
					fmt.Sprintf("The 'request-id' header has not founded [%s]",
						fmt.Sprintf("%s %s",
							html.EscapeString(r.Method),
							html.EscapeString(r.URL.Path))))
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			}
			if toCheck {
				if len(requestID) != 29 {
					log.Printf("ERROR [%s: %s]", "FXZ5Jc1t",
						fmt.Sprintf("The 'request-id' header has an invalid value [%s]",
							fmt.Sprintf("%s %s %s",
								html.EscapeString(r.Method),
								html.EscapeString(r.URL.Path),
								html.EscapeString(r.Header.Get("request-id")))))
					http.Error(w, "400 Bad Request", http.StatusBadRequest)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
