/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package rest

import (
	"github.com/dmalix/financelime-authorization/app/information/service"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersion(t *testing.T) {

	service.MockData.Expected.Error = nil

	request, err := http.NewRequest(
		http.MethodGet,
		"/",
		nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	logger := new(zap.Logger)
	newContextGetter := new(middleware.MockDescription)
	newService := new(service.MockDescription)
	newAPI := NewREST(newContextGetter, newService)
	handler := newAPI.Version(logger)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
