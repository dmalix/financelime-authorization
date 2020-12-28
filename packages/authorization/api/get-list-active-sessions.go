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
			errorMessage      string
			encryptedUserData []byte
			sessions          []models.Session
		)

		if r.Context().Value(ContextEncryptedUserData) == nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get Context from the request [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("500 Internal Server Error"), http.StatusBadRequest)
			return
		}

		encryptedUserData = r.Context().Value(ContextEncryptedUserData).([]byte)

		sessions, err = h.service.GetListActiveSessions(encryptedUserData)
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		responseBody, err = json.Marshal(sessions)
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusOK
		contentType = ContentTypeApplicationJson

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
