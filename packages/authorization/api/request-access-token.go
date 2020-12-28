/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/utils/responder"
	"github.com/dmalix/financelime-authorization/utils/trace"
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

		remoteAddr = r.Header.Get("X-Real-IP")
		if len(remoteAddr) == 0 {
			remoteAddr = r.RemoteAddr
		}

		response.PublicSessionID, response.AccessJWT, response.RefreshJWT, err =
			h.service.RequestAccessToken(props.Email, props.Password, props.ClientID, remoteAddr, props.Device)
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to request an access token"
			switch domainErrorCode {
			case "PROPS": // One or more of the input parameters are invalid
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			case "USER_NOT_FOUND": // User is not found
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
				return
			}
		}

		responseBody, err = json.Marshal(response)
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusOK
		contentType = "application/json;charset=utf-8"

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
