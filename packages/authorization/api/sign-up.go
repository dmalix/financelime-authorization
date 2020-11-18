/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/responder"
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
			errLabel        string
			domainErrorCode string
			errorMessage    string
		)

		if strings.ToLower(r.Header.Get("content-type")) != "application/json;charset=utf-8" {
			errLabel = "26Pi82rl"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			errLabel = "w5a7C38O"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = r.Body.Close()
		if err != nil {
			errLabel = "8w5a7C3O"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to close a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			errLabel = "jlgeF0it"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to convert a body props to struct [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
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
				errLabel = "jInpoLV5"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			case "USER_ALREADY_EXIST": // a user with the email you specified already exists
				errLabel = "5Ig7X4Sv"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			case "INVITE_NOT_EXIST_EXPIRED": // the invite code does not exist or is expired
				errLabel = "61H2IR2f"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			case "INVITE_LIMIT": // the limit for issuing this invite code has been exhausted
				errLabel = "pZ4fgc9k"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			default:
				errLabel = "e3YlkJHc"
				log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
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
