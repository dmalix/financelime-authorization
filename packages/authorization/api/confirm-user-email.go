/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"github.com/dmalix/financelime-rest-api/utils/responder"
	"github.com/dmalix/financelime-rest-api/utils/url"
	"log"
	"net/http"
	"strings"
)

//	Confirm user email
func (h *Handler) ConfirmUserEmail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			confirmationKey string
			message         string
			responseBody    []byte
			statusCode      int
			contentType     string
			err             error
			errLabel        string
			domainErrorCode string
			errorMessage    string
		)

		confirmationKey, err = url.GetPathValue(r.URL.Path, 1)
		if err != nil {
			errLabel = "LVjInpo5"
			log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		message, err = h.service.ConfirmUserEmail(confirmationKey)
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to confirm user email"
			switch domainErrorCode {
			case "CONFIRMATION_KEY_NOT_VALID": // the confirmation key not valid
				errLabel = "C2V0NqJm"
				log.Printf("ERROR [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			default:
				errLabel = "GBQV0Zc1"
				log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
				w.Header().Add("error-label", errLabel)
				http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
				return
			}
		}

		statusCode = http.StatusOK
		responseBody = []byte(message)
		contentType = "text/plain;charset=utf-8"

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
