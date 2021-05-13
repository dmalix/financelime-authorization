package middleware

import (
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"html"
	"log"
	"net/http"
)

func (middleware *Middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const (
			isRequired = true
			toCheck    = true
		)

		if isRequired {
			requestID := r.Header.Get("request-id")
			if len(requestID) == 0 {
				log.Printf("%s: %s %s", "ERROR", trace.GetCurrentPoint(),
					fmt.Sprintf("The 'request-id' header not found [%s]",
						fmt.Sprintf("%s %s",
							html.EscapeString(r.Method),
							html.EscapeString(r.URL.Path))))
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			}
			if toCheck {
				if len(requestID) != 29 {
					log.Printf("%s: %s %s", "ERROR", trace.GetCurrentPoint(),
						fmt.Sprintf("The 'request-id' header is invalid [%s]",
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
