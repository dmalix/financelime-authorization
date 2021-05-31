package middleware

import (
	"context"
	"github.com/dmalix/financelime-authorization/packages/jwt"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handler200() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		//noinspection GoUnhandledErrorResult
		w.Write(nil)
	})
}

func TestRequestID_400(t *testing.T) {

	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(request.Context(), ContextKeyRemoteAddr, "127.0.0.1")

	authorizationAPIMiddlewareConfig := ConfigMiddleware{
		RequestIDRequired: true,
		RequestIDCheck:    true,
	}

	jwtManager := jwt.NewToken(
		"12345",
		jwt.ParamSigningAlgorithmHS256,
		"issuer",
		"subject",
		1000,
		1000)

	authorizationAPIMiddleware := NewMiddleware(
		authorizationAPIMiddlewareConfig,
		jwtManager)

	request.Header.Add("request-id", "1234")
	responseRecorder := httptest.NewRecorder()
	logger, _ := zap.NewProduction()
	mwFunc := authorizationAPIMiddleware.RequestID(logger)
	handler := mwFunc(handler200())
	handler.ServeHTTP(responseRecorder, request.WithContext(ctx))
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestRequestID_200(t *testing.T) {

	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(request.Context(), ContextKeyRemoteAddr, "127.0.0.1")

	authorizationAPIMiddlewareConfig := ConfigMiddleware{
		RequestIDRequired: true,
		RequestIDCheck:    true,
	}

	jwtManager := jwt.NewToken(
		"12345",
		jwt.ParamSigningAlgorithmHS256,
		"issuer",
		"subject",
		1000,
		1000)

	authorizationAPIMiddleware := NewMiddleware(
		authorizationAPIMiddlewareConfig,
		jwtManager)

	request.Header.Add("request-id", "abcda12b12c12d12")
	responseRecorder := httptest.NewRecorder()
	logger, _ := zap.NewProduction()
	mwFunc := authorizationAPIMiddleware.RequestID(logger)
	handler := mwFunc(handler200())
	handler.ServeHTTP(responseRecorder, request.WithContext(ctx))
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
