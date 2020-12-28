/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/responder"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//	Create new user
func (h *Handler) SignUp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type incomingProps struct {
			Email      string `json:"email"`
			InviteCode string `json:"inviteCode"`
			Language   string `json:"language"`
		}

		var (
			props           incomingProps
			body            []byte
			responseBody    []byte
			statusCode      int
			contentType     string
			remoteAdrr      string
			err             error
			domainErrorCode string
			errorMessage    string
		)

		if strings.ToLower(r.Header.Get("content-type")) != "application/json;charset=utf-8" {
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			http.Error(w, fmt.Sprintf("400 Bad Request"), http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("400 Bad Request"), http.StatusBadRequest)
			return
		}
		err = r.Body.Close()
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to close a body [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("500 Internal Server Error"), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to convert a body props to struct [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("400 Bad Request"), http.StatusBadRequest)
			return
		}

		remoteAdrr = r.Header.Get("X-Real-IP")
		if len(remoteAdrr) == 0 {
			remoteAdrr = r.RemoteAddr
		}

		err = h.service.SignUp(props.Email, props.Language, props.InviteCode, remoteAdrr)
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to Sign Up"
			switch domainErrorCode {
			case "PROPS": // one or more of the input parameters are invalid
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			case "USER_ALREADY_EXIST": // a user with the email you specified already exists
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			case "INVITE_NOT_EXIST_EXPIRED": // the invite code does not exist or is expired
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			case "INVITE_LIMIT": // the limit for issuing this invite code has been exhausted
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
				return
			}
		}

		statusCode = http.StatusNoContent
		responseBody = nil
		contentType = ""

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
