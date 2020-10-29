/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
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

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_a1_ServiceError(t *testing.T) {

	service.Mock.Values.SignUp.ExpectedError = errors.New(fmt.Sprintf("%s:", "a1"))

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
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_409_b1_ServiceError(t *testing.T) {

	service.Mock.Values.SignUp.ExpectedError = errors.New(fmt.Sprintf("%s:", "b1"))

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
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_409_b2_ServiceError(t *testing.T) {

	service.Mock.Values.SignUp.ExpectedError = errors.New(fmt.Sprintf("%s:", "b2"))

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
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_409_b3_ServiceError(t *testing.T) {

	service.Mock.Values.SignUp.ExpectedError = errors.New(fmt.Sprintf("%s:", "b3"))

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
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_500_DEFAULT_ServiceError(t *testing.T) {

	service.Mock.Values.SignUp.ExpectedError = errors.New("ServerError")

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
	request.Header.Add("content-type", "application/json;charset=utf-8")

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSignUp_202(t *testing.T) {

	service.Mock.Values.SignUp.Props.Email = "testuser@financelime.com"
	service.Mock.Values.SignUp.Props.InviteCode = "testInviteCode"
	service.Mock.Values.SignUp.Props.Language = "en"
	service.Mock.Values.SignUp.Props.RemoteAddr = "127.0.0.1"

	service.Mock.Values.SignUp.ExpectedError = nil

	props := map[string]interface{}{
		"email":      service.Mock.Values.SignUp.Props.Email,
		"inviteCode": service.Mock.Values.SignUp.Props.InviteCode,
		"language":   service.Mock.Values.SignUp.Props.Language,
		"remoteAddr": service.Mock.Values.SignUp.Props.RemoteAddr,
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
	request.Header.Add("X-Real-IP", service.Mock.Values.SignUp.Props.RemoteAddr)

	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockType)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}
}
