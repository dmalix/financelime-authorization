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

func TestRequestUserPasswordReset(t *testing.T) {

	service.MockData.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "/", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.RequestUserPasswordReset()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}
