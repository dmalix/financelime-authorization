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

//	Request an access token
func (h *Handler) RefreshAccessToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type incomingProps struct {
			RefreshToken string `json:"refreshToken"`
		}

		type outgoingResponse struct {
			PublicSessionID string `json:"sessionID"`
			AccessJWT       string `json:"accessToken"`
			RefreshJWT      string `json:"refreshToken"`
		}

		var (
			props           incomingProps
			body            []byte
			response        outgoingResponse
			responseBody    []byte
			statusCode      int
			contentType     string
			remoteAddr      string
			err             error
			errLabel        string
			domainErrorCode string
			errorMessage    string
		)

		if strings.ToLower(r.Header.Get("content-type")) != "application/json;charset=utf-8" {
			errLabel = "oC7ohCie"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			errLabel = "aiGee8Ei"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = r.Body.Close()
		if err != nil {
			errLabel = "Jiek7ooM"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to close a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			errLabel = "aiBein9e"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to convert a body props to struct [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}

		remoteAddr = r.Header.Get("X-Real-IP")
		if len(remoteAddr) == 0 {
			remoteAddr = r.RemoteAddr
		}

		response.PublicSessionID, response.AccessJWT, response.RefreshJWT, err =
			h.service.RefreshAccessToken(props.RefreshToken, remoteAddr)
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to request an access token"
			switch domainErrorCode {
			case "INVALID_REFRESH_TOKEN", "USER_NOT_FOUND": // One or more of the input parameters are invalid
				errLabel = "RtueOfg6"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			default:
				errLabel = "wCnUTP27"
				log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
				return
			}
		}

		responseBody, err = json.Marshal(response)
		if err != nil {
			errLabel = "vWu8wbXv"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusOK
		contentType = "application/json;charset=utf-8"

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
