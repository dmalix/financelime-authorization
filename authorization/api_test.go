/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPISignUp(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	request, err := http.NewRequest(http.MethodPost, "/v1/signup", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", contentTypeApplicationJson)
	request.Header.Add("X-Real-IP", ServiceMockData.Props.RemoteAddr)

	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.signUp(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIConfirmUserEmail(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	propsByte, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodGet, "/v1/u/key", bytes.NewBuffer(propsByte))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.confirmUserEmail(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPICreateAccessToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/v1/oauth/token", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", contentTypeApplicationJson)

	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.createAccessToken(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIRefreshAccessToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodPut, "/v1/oauth/token", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("content-type", contentTypeApplicationJson)

	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.refreshAccessToken(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIRevokeRefreshToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{
		"sessionID": "870bd06be766720b7348f6baf946355b71d23401978f7199b8437f52377f62e1",
	}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodDelete, "/v1/oauth/sessions", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add("content-type", contentTypeApplicationJson)

	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.revokeRefreshToken(ctx)

	rctx := request.Context()
	rctx = context.WithValue(rctx, middleware.ContextPublicSessionID, "PublicSessionID")
	rctx = context.WithValue(rctx, middleware.ContextEncryptedUserData, []byte("EncryptedUserData"))

	request = request.WithContext(rctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIRequestUserPasswordReset(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ServiceMockData.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/resetpassword", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add("content-type", contentTypeApplicationJson)

	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.requestUserPasswordReset(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIGetListActiveSessions(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ServiceMockData.Expected.Error = nil

	request, err := http.NewRequest(http.MethodGet, "/v1/oauth/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	service := new(ServiceMockDescription)
	api := NewAPI(service)
	handler := api.getListActiveSessions(ctx)

	rctx := request.Context()
	rctx = context.WithValue(rctx, middleware.ContextEncryptedUserData, []byte("test_data"))

	request = request.WithContext(rctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
