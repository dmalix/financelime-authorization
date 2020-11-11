/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/utils/responder"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//	Request an access token
func (h *Handler) RequestAccessToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type incomingProps struct {
			Email    string        `json:"email"`
			Password string        `json:"password"`
			ClientID string        `json:"client_id"`
			Device   models.Device `json:"device"`
		}

		type outgoingResponse struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}

		var (
			props           incomingProps
			body            []byte
			response        outgoingResponse
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
			errLabel = "Gwn3ryea"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			errLabel = "N4lidcri"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = r.Body.Close()
		if err != nil {
			errLabel = "5UMtv0YJ"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to close a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			errLabel = "TALDtv9L"
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

		response.AccessToken, response.RefreshToken, err =
			h.service.RequestAccessToken(props.Email, props.Password, props.ClientID, remoteAdrr, props.Device)
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to request an access token"
			switch domainErrorCode {
			case "PROPS": // One or more of the input parameters are invalid
				errLabel = "e0c1Dbcq"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			case "USER_NOT_FOUND": // User is not found
				errLabel = "tueOfg6R"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "409 Conflict", http.StatusConflict)
				return
			default:
				errLabel = "CnUTwP27"
				log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
				return
			}
		}

		responseBody, err = json.Marshal(response)
		if err != nil {
			errLabel = "Wu8wbXvv"
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
