/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dmalix/financelime-authorization/app/authorization/service"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPISignUp(t *testing.T) {

	authService := new(service.Mock)

	authService.Props.Email = "user@domain.com"
	authService.Props.InviteCode = "invite_code"
	authService.Props.Language = "abc"
	authService.Props.RemoteAddr = "127.0.0.1"

	authService.Expected.Error = nil

	props := map[string]interface{}{
		"email":      authService.Props.Email,
		"inviteCode": authService.Props.InviteCode,
		"language":   authService.Props.Language,
		"remoteAddr": authService.Props.RemoteAddr,
	}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add(headerKeyContentType, headerValueApplicationJson)
	request.Header.Add("X-Real-IP", authService.Props.RemoteAddr)

	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.SignUp(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIConfirmUserEmail(t *testing.T) {

	authService := new(service.Mock)

	authService.Expected.Error = nil

	props := map[string]interface{}{}

	propsByte, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "", bytes.NewBuffer(propsByte))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.ConfirmUserEmail(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPICreateAccessToken(t *testing.T) {

	authService := new(service.Mock)

	authService.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add(headerKeyContentType, headerValueApplicationJson)

	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.CreateAccessToken(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIRefreshAccessToken(t *testing.T) {

	authService := new(service.Mock)

	authService.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add(headerKeyContentType, headerValueApplicationJson)

	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.RefreshAccessToken(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIRevokeRefreshToken(t *testing.T) {

	authService := new(service.Mock)

	authService.Expected.Error = nil

	props := map[string]interface{}{
		"sessionID": "870bd06be766720b7348f6baf946355b71d23401978f7199b8437f52377f62e1",
	}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add(headerKeyContentType, headerValueApplicationJson)

	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.RevokeRefreshToken(logger)

	rctx := request.Context()
	rctx = context.WithValue(rctx, middleware.ContextKeyPublicSessionID, "PublicSessionID")
	rctx = context.WithValue(rctx, middleware.ContextKeyEncryptedJWTData, []byte("EncryptedJWTData"))

	request = request.WithContext(rctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIRequestUserPasswordReset(t *testing.T) {

	authService := new(service.Mock)

	authService.Expected.Error = nil

	props := map[string]interface{}{}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("", "", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add(headerKeyContentType, headerValueApplicationJson)

	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.RequestUserPasswordReset(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPIGetListActiveSessions(t *testing.T) {

	authService := new(service.Mock)

	authService.Expected.Error = nil

	request, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	logger, _ := zap.NewProduction()
	contextGetter := new(middleware.MockDescription)

	authREST := NewREST(contextGetter, authService)
	handler := authREST.GetListActiveSessions(logger)

	rctx := request.Context()
	rctx = context.WithValue(rctx, middleware.ContextKeyEncryptedJWTData, []byte("test_data"))

	request = request.WithContext(rctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
