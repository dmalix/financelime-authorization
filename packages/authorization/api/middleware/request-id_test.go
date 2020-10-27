package middleware

import (
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

	request.Header.Add("request-id", "1234")
	responseRecorder := httptest.NewRecorder()
	handler := RequestID(handler200())
	handler.ServeHTTP(responseRecorder, request)
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

	request.Header.Add("request-id", "K7800-H7625-Z5852-N1693-K1972")
	responseRecorder := httptest.NewRecorder()
	handler := RequestID(handler200())
	handler.ServeHTTP(responseRecorder, request)
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
