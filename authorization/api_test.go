/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPISignUp(t *testing.T) {

	ServiceMockData.Props.Email = "user@domain.com"
	ServiceMockData.Props.InviteCode = "invite_code"
	ServiceMockData.Props.Language = "abc"
	ServiceMockData.Props.RemoteAddr = "127.0.0.1"

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{
		"email":      ServiceMockData.Props.Email,
		"inviteCode": ServiceMockData.Props.InviteCode,
		"language":   ServiceMockData.Props.Language,
		"remoteAddr": ServiceMockData.Props.RemoteAddr,
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
	request.Header.Add("X-Real-IP", ServiceMockData.Props.RemoteAddr)

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.signUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIConfirmUserEmail_200(t *testing.T) {

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	propsByte, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "/acue/12345", bytes.NewBuffer(propsByte))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.confirmUserEmail()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIConfirmUserEmail_500(t *testing.T) {

	ServiceMockData.Expected.Error = errors.New("SERVICE_ERROR")

	props := map[string]interface{}{}

	propsByte, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "/acue/12345", bytes.NewBuffer(propsByte))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.confirmUserEmail()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestAPIConfirmUserEmail_404__CONFIRMATION_KEY_NOT_FOUND(t *testing.T) {

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	propsByte, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "/acue", bytes.NewBuffer(propsByte))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.confirmUserEmail()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestAPIConfirmUserEmail_404__CONFIRMATION_KEY_NOT_VALID(t *testing.T) {

	errorMessage := "failed to confirm user email"
	errLabel := "oLjInpV5"
	err := errors.New("SERVICE_ERROR")
	ServiceMockData.Expected.Error = errors.New(fmt.Sprintf("ERROR [%s:%s[%v]]", errLabel, errorMessage, err))

	props := map[string]interface{}{}

	propsByte, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "/acue", bytes.NewBuffer(propsByte))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.confirmUserEmail()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestAPIRequestAccessToken(t *testing.T) {

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(
		"",
		"/",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.createAccessToken()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIRefreshAccessToken(t *testing.T) {

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(
		"",
		"/",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.refreshAccessToken()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIRevokeRefreshToken(t *testing.T) {

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{
		"sessionID": "870bd06be766720b7348f6baf946355b71d23401978f7199b8437f52377f62e1",
	}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(
		http.MethodDelete,
		"/v1/oauth/sessions",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.revokeRefreshToken()

	ctx := request.Context()
	ctx = context.WithValue(ctx, contextPublicSessionID, "PublicSessionID")
	ctx = context.WithValue(ctx, contextEncryptedUserData, []byte("EncryptedUserData"))

	request = request.WithContext(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIRequestUserPasswordReset(t *testing.T) {

	ServiceMockData.Expected.Error = nil

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

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.requestUserPasswordReset()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIGetListActiveSessions(t *testing.T) {

	ServiceMockData.Expected.Error = nil

	request, err := http.NewRequest(
		http.MethodGet,
		"/v1/oauth/sessions",
		nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	authService := new(ServiceMockDescription)
	newHandler := NewAPI(authService)
	handler := newHandler.getListActiveSessions()

	ctx := request.Context()
	ctx = context.WithValue(ctx,
		contextEncryptedUserData,
		[]byte("test_data"))

	request = request.WithContext(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
