/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package middleware

import (
	"context"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/packages/authorization/api"
	"html"
	"log"
	"net/http"
	"strings"
)

func (m *Middleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			errLabel      string
			err           error
			authorization string
			jwtTokenArr   []string
			jwtData       models.JwtData
			ctx           context.Context
		)

		// Get an authorization token from the header

		authorization = r.Header.Get("authorization")
		if len(authorization) == 0 {
			errLabel = "NUV9WZfq"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("The 'authorization' header has not founded [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("401 Unauthorized [%s]", errLabel), http.StatusUnauthorized)
			return
		}

		// Validate Token and data extract for identification

		jwtTokenArr = strings.Split(strings.TrimSpace(html.EscapeString(authorization)), " ")
		if len(jwtTokenArr) != 2 {
			errLabel = "9ZNUVWfq"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("The 'authorization' header has not founded [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("401 Unauthorized [%s]", errLabel), http.StatusUnauthorized)
			return
		}

		if strings.ToLower(jwtTokenArr[0]) != "bearer" {
			errLabel = "u5FCvHqa"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("Invalid JWT-Token [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("401 Unauthorized [%s]", errLabel), http.StatusUnauthorized)
			return
		}

		jwtData, err = m.jwt.VerifyToken(jwtTokenArr[1])
		if err != nil {
			errLabel = "hk7LCW2T"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Invalid JWT-Token [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("401 Unauthorized [%s]", errLabel), http.StatusUnauthorized)
			return
		}

		if jwtData.Payload.Purpose != "access" {
			errLabel = "2vRqSqVa"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("Invalid JWT-Token [%s]",
					fmt.Sprintf("%s %s",
						html.EscapeString(r.Method),
						html.EscapeString(r.URL.Path))))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("401 Unauthorized [%s]", errLabel), http.StatusUnauthorized)
			return
		}

		ctx = r.Context()
		ctx = context.WithValue(ctx, api.ContextPublicSessionID, jwtData.Payload.PublicSessionID)
		ctx = context.WithValue(ctx, api.ContextEncryptedUserData, jwtData.Payload.UserData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
