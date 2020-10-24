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
			label        string
		)

		if strings.ToLower(r.Header.Get("content-type")) != "application/json;charset=utf-8" {
			label = "26Pi82rl"
			log.Printf("ERROR [%s: %s]", label,
				fmt.Sprintf("Header 'content-type:application/json;charset=utf-8' not found [%s]",
					responder.Message(r)))
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", label), http.StatusBadRequest)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			label = "w5a7C38O"
			log.Printf("ERROR [%s: %s [%s]]", label,
				fmt.Sprintf("Failed to get a body [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", label), http.StatusBadRequest)
			return
		}
		err = r.Body.Close()
		if err != nil {
			label = "8w5a7C3O"
			log.Printf("ERROR [%s: %s [%s]]", label,
				fmt.Sprintf("Failed to close a body [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("500 Internal Server Error [%s]", label), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &props)
		if err != nil {
			label = "jlgeF0it"
			log.Printf("ERROR [%s: %s [%s]]", label,
				fmt.Sprintf("Failed to convert a body props to struct [%s]", responder.Message(r)),
				err)
			http.Error(w, fmt.Sprintf("400 Bad Request [%s]", label), http.StatusBadRequest)
			return
		}

		remoteAdrr = r.Header.Get("X-Real-IP")
		if len(remoteAdrr) == 0 {
			remoteAdrr = r.RemoteAddr
		}

		err = h.service.SignUp(props.Email, props.InviteCode, props.Language, remoteAdrr)
		if err != nil {
			label = "G0bFFCuq"
			log.Printf("FATAL [%s: Failed to Sign Up: [%s]]", label, err)
			http.Error(w, fmt.Sprintf("500 Server Internal Error [%s]", label), http.StatusInternalServerError)
			return
		}

		statusCode = 202
		responseBody = nil
		contentType = ""

		responder.Response(w, r, responseBody, statusCode, contentType)
		return
	})
}
