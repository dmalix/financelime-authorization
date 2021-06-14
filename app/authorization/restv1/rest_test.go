package restv1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dmalix/authorization-service/app/authorization/model"
	"github.com/dmalix/authorization-service/app/authorization/service"
	"github.com/dmalix/jwt"
	"github.com/dmalix/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPISignUp1(t *testing.T) {

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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
	handler := authREST.SignUpStep1(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestAPISignUp2(t *testing.T) {

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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
	handler := authREST.SignUpStep2(logger)

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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
	handler := authREST.RevokeRefreshToken(logger)

	rctx := request.Context()
	rctx = context.WithValue(rctx, middleware.ContextKeyJwt, jwt.Token{})

	request = request.WithContext(rctx)

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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
	handler := authREST.GetListActiveSessions(logger)

	rctx := request.Context()
	rctx = context.WithValue(rctx, middleware.ContextKeyJwt, jwt.Token{})

	request = request.WithContext(rctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAPIResetUserPassword1(t *testing.T) {

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

	authREST := NewRest(model.ConfigRest{}, contextGetter, authService)
	handler := authREST.ResetUserPasswordStep1(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}
