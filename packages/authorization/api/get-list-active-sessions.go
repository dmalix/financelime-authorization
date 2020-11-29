/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/utils/responder"
	"log"
	"net/http"
)

//	Create new user
func (h *Handler) GetListActiveSessions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			responseBody      []byte
			statusCode        int
			contentType       string
			err               error
			errLabel          string
			errorMessage      string
			encryptedUserData []byte
			sessions          []models.Session
		)

		if r.Context().Value(ContextEncryptedUserData) == nil {
			errLabel = "a7wC583O"
			log.Printf("ERROR [%s: %s [%s]]", errLabel,
				fmt.Sprintf("Failed to get Context from the request [%s]", responder.Message(r)),
				err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", errLabel), http.StatusBadRequest)
			return
		}

		encryptedUserData = r.Context().Value(ContextEncryptedUserData).([]byte)

		sessions, err = h.service.GetListActiveSessions(encryptedUserData)
		if err != nil {
			errLabel = "8XvWuwbv"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		responseBody, err = json.Marshal(sessions)
		if err != nil {
			errLabel = "XvWu8wbv"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusOK
		contentType = ContentTypeApplicationJson

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
