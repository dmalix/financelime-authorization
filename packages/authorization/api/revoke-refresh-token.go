/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/utils/responder"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//	Create new user
func (h *Handler) RevokeRefreshToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type incomingProps struct {
			PublicSessionID string `json:"sessionID"`
		}

		var (
			props             incomingProps
			responseBody      []byte
			statusCode        int
			contentType       string
			err               error
			errLabel          string
			errorMessage      string
			body              []byte
			publicSessionID   string
			encryptedUserData []byte
			sessions          []models.Session
		)

		if strings.ToLower(r.Header.Get("content-type")) != "application/json;charset=utf-8" {
			errLabel = "826Pi2rl"
			log.Printf("ERROR [%s: %s]", errLabel,
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			errLabel = "8w5a7C3O"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = r.Body.Close()
		if err != nil {
			errLabel = "7C8w5a3O"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to close a body [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			errLabel = "F0jlgeit"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to convert a body props to struct [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", errLabel), http.StatusBadRequest)
			return
		}

		if r.Context().Value(ContextEncryptedUserData) == nil {
			errLabel = "wC83a75O"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to get ContextEncryptedUserData from the request [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
			return
		}

		encryptedUserData = r.Context().Value(ContextEncryptedUserData).([]byte)

		if len(props.PublicSessionID) == 0 {
			if r.Context().Value(ContextPublicSessionID) == nil {
				errLabel = "83a75wCO"
				log.Printf("ERROR [%s: %s [%s]]", errLabel,
					fmt.Sprintf("Failed to get ContextPublicSessionID from the request [%s]", responder.Message(r)),
					err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
				return
			}
			publicSessionID = r.Context().Value(ContextPublicSessionID).(string)
		} else {
			publicSessionID = props.PublicSessionID
		}

		err = h.service.RevokeRefreshToken(encryptedUserData, publicSessionID)
		if err != nil {
			errLabel = "vWbuw8Xv"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		responseBody, err = json.Marshal(sessions)
		if err != nil {
			errLabel = "8wbXvWuv"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusNoContent
		contentType = ContentTypeTextPlain

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}