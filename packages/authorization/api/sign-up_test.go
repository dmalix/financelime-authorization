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

	authService := new(service.MockDescription)
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

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_400_ServiceError__PROPS(t *testing.T) {

	service.MockData.Expected.Error = errors.New(fmt.Sprintf("%s:", "PROPS"))

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

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSignUp_409_ServiceError__USER_ALREADY_EXIST(t *testing.T) {

	service.MockData.Expected.Error = errors.New(fmt.Sprintf("%s:", "USER_ALREADY_EXIST"))

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

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_409_ServiceError__INVITE_NOT_EXIST_EXPIRED(t *testing.T) {

	service.MockData.Expected.Error = errors.New(fmt.Sprintf("%s:", "INVITE_NOT_EXIST_EXPIRED"))

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

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_409_ServiceError__INVITE_LIMIT(t *testing.T) {

	service.MockData.Expected.Error = errors.New(fmt.Sprintf("%s:", "INVITE_LIMIT"))

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

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestSignUp_500_ServiceError__SYSTEM(t *testing.T) {

	service.MockData.Expected.Error = errors.New("ServiceError")

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

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.SignUp()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSignUp_202(t *testing.T) {

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

	if status := responseRecorder.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}
}
