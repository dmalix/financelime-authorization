/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"bytes"
	"encoding/json"
	"github.com/dmalix/financelime-rest-api/packages/authorization/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp_400_NoHeaderContentType(t *testing.T) {

	props := map[string]interface{}{}

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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_InvalidHeaderContentType(t *testing.T) {

	props := map[string]interface{}{}

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

	request.Header.Add("content-type", "1234")

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_InvalidEmailParam(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "testInviteCode"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      "1234",
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_InvalidInviteCodeParam(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "testInviteCode"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": "1234",
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_InvalidLanguageParam(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "testInviteCode"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   "1234",
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_InvalidRemoteAddrParam(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "testInviteCode"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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
	request.Header.Add("X-Real-IP", "ParamRemoteAddrIsNotValid")

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_409_TheEmailExists(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "email.exists@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "testInviteCode"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_409_TheInviteCodeDoesNotExistOrIsExpired(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "InviteCodeErrorFL104"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_409_TheLimitForIssuingThisInviteCodeHasBeenExhausted(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "InviteCodeErrorFL105"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_202(t *testing.T) {

	service.ServiceMockValue.Props.SignUp.Email = "testuser@financelime.com"
	service.ServiceMockValue.Props.SignUp.InviteCode = "testInviteCode"
	service.ServiceMockValue.Props.SignUp.Language = "en"
	service.ServiceMockValue.Props.SignUp.RemoteAddr = "127.0.0.1"

	props := map[string]interface{}{
		"email":      service.ServiceMockValue.Props.SignUp.Email,
		"inviteCode": service.ServiceMockValue.Props.SignUp.InviteCode,
		"language":   service.ServiceMockValue.Props.SignUp.Language,
		"remoteAddr": service.ServiceMockValue.Props.SignUp.RemoteAddr,
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

	responseRecorder := httptest.NewRecorder()

	serviceMock := new(service.AuthorizationServiceMock)
	newHandler := NewHandler(serviceMock)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}
}

