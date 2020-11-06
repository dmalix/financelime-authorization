/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/utils/responder"
	"log"
	"net/http"
)

func (handler *Handler) Dist() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			responseBody []byte
			statusCode   int
			contentType  string
			err          error
			errLabel     string
			errorMessage string
			dist         models.Dist
			distJSON     []byte
		)

		dist.Version, dist.Build, err = handler.service.Dist()
		if err != nil {
			errorMessage = "failed to get Dist"
			errLabel = "3YlJekHc"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		distJSON, err = json.Marshal(&dist)
		if err != nil {
			errorMessage = "failed to convert the dist data to JSON-format"
			errLabel = "Je3YlkHc"
			log.Printf("FATAL [%s:%s[%s]]", errLabel, errorMessage, err)
			w.Header().Add("error-label", errLabel)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusOK
		responseBody = distJSON
		contentType = ""

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
