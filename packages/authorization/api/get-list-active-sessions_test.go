/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"context"
	"github.com/dmalix/financelime-authorization/packages/authorization/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetListActiveSessions(t *testing.T) {

	service.MockData.Expected.Error = nil

	request, err := http.NewRequest(
		http.MethodGet,
		"/v1/oauth/sessions",
		nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	authService := new(service.MockDescription)
	newHandler := NewHandler(authService)
	handler := newHandler.GetListActiveSessions()

	ctx := request.Context()
	ctx = context.WithValue(ctx,
		ContextEncryptedUserData,
		[]byte("test_data"))

	request = request.WithContext(ctx)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
