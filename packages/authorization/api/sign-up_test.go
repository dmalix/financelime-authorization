/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"bytes"
	"encoding/json"
	"github.com/dmalix/financelime-authorization/packages/authorization/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {

	service.MockData.Props.Email = "user@domain.com"
	service.MockData.Props.InviteCode = "invite_code"
	service.MockData.Props.Language = "abc"
	service.MockData.Props.RemoteAddr = "127.0.0.1"

	service.MockData.Expected.Error = nil

	props := map[string]interface{}{
		"email":      service.MockData.Props.Email,
		"inviteCode": service.MockData.Props.InviteCode,
		"language":   service.MockData.Props.Language,
		"remoteAddr": service.MockData.Props.RemoteAddr,
	}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(
		"POST",
		"/authorization/signup",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", "application/json;charset=utf-8")
	request.Header.Add("X-Real-IP", service.MockData.Props.RemoteAddr)

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}
