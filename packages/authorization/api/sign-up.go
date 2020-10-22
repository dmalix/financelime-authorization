/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-rest-api/utils/responder"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//	Create new account
func (h *Handler) SignUp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type incomingProps struct {
			Email      string `json:"email"`
			InviteCode string `json:"inviteCode"`
			Language   string `json:"language"`
		}

		var (
			props        incomingProps
			err          error
			body         []byte
			responseBody []byte
			statusCode   int
			contentType  string
			remoteAdrr   string
		)

		if strings.ToLower(r.Header.Get("content-type")) != "application/json;charset=utf-8" {
			log.Printf("ERROR [%s: %s]", "26Pi82rl",
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", "w5a7C38O",
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", "jlgeF0it",
				fmt.Sprintf("Failed to convert a body props to struct [%s]", responder.Message(r)),
				err)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		remoteAdrr = r.Header.Get("X-Real-IP")
		if len(remoteAdrr) == 0 {
			remoteAdrr = r.RemoteAddr
		}

		err = h.service.SignUp(props.Email, props.InviteCode, props.Language, remoteAdrr)
		if err != nil {
			log.Printf("FATAL [%s: Failed to Sign Up: [%s]]", "G0bFFCuq", err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		statusCode = 202
		responseBody = nil
		contentType = ""

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
