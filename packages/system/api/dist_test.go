/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"github.com/dmalix/financelime-rest-api/packages/system/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDist(t *testing.T) {

	service.Dist_MockData.Expected.Error = nil

	request, err := http.NewRequest(
		http.MethodGet,
		"/",
		nil)
	if err != nil {
		t.Fatal(err)
	}
	responseRecorder := httptest.NewRecorder()

	newService := new(service.Dist_MockDescription)
	newHandler := NewHandler(newService)
	handler := newHandler.Dist()

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
