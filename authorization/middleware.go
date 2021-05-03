/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"context"
	"fmt"
	jwt2 "github.com/dmalix/financelime-authorization/packages/jwt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"html"
	"log"
	"net/http"
	"strings"
)

type ConfigMiddleware struct {
	RequestIDRequired bool
	RequestIDCheck    bool
}

type Middleware struct {
	config ConfigMiddleware
	jwt    jwt2.Jwt
}

func NewMiddleware(
	config ConfigMiddleware,
	jwt jwt2.Jwt) *Middleware {
	return &Middleware{
		config: config,
		jwt:    jwt,
	}
}

func (m *Middleware) authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			err           error
			authorization string
			jwtTokenArr   []string
			jwtData       jwt2.JwtData
			ctx           context.Context
		)

		// Get an authorization token from the header

		authorization = r.Header.Get("authorization")
		if len(authorization) == 0 {
			log.Printf("%s: %s %s", "ERROR", trace.GetCurrentPoint(),
				fmt.Sprintf("The 'authorization' header has not founded [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate Token and data extract for identification

		jwtTokenArr = strings.Split(strings.TrimSpace(html.EscapeString(authorization)), " ")
		if len(jwtTokenArr) != 2 {
			log.Printf("%s: %s %s", "ERROR", trace.GetCurrentPoint(),
				fmt.Sprintf("The 'authorization' header has not founded [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		if strings.ToLower(jwtTokenArr[0]) != "bearer" {
			log.Printf("%s: %s %s", "ERROR", trace.GetCurrentPoint(),
				fmt.Sprintf("Invalid JWT-Token [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		jwtData, err = m.jwt.VerifyToken(jwtTokenArr[1])
		if err != nil {
			log.Printf("%s: %s %s [%s]", "ERROR", trace.GetCurrentPoint(),
				fmt.Sprintf("Invalid JWT-Token [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))),
				err)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		if jwtData.Payload.Purpose != "access" {
			log.Printf("%s: %s %s", "ERROR", trace.GetCurrentPoint(),
				fmt.Sprintf("Invalid JWT-Token [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			http.Error(w, "401 Unauthorized [%s]", http.StatusUnauthorized)
			return
		}

		ctx = r.Context()
		ctx = context.WithValue(ctx, contextPublicSessionID, jwtData.Payload.PublicSessionID)
		ctx = context.WithValue(ctx, contextEncryptedUserData, jwtData.Payload.UserData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) requestID(next http.Handler) http.Handler {
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
