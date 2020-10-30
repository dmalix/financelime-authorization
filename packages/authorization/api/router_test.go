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

func TestRouter(t *testing.T) {

	var (
		res *http.Response
		err error
	)

	mux := http.NewServeMux()
	authService := new(service.SignUp_MockDescription)
	Router(mux, authService)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// StatusNotFound

	res, err = http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code for /not-exists is wrong. Have: %d, want: %d.",
			res.StatusCode, http.StatusNotFound)
	}

	// StatusMethodNotAllowed

	res, err = http.Get(ts.URL + "/authorization/signup")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status code for /authorization/signup is wrong. Have: %d, want: %d.",
			res.StatusCode, http.StatusMethodNotAllowed)
	}

	// StatusAccepted

	service.SignUp_MockData.Props.Email = "testuser@financelime.com"
	service.SignUp_MockData.Props.InviteCode = "testInviteCode"
	service.SignUp_MockData.Props.Language = "en"
	service.SignUp_MockData.Props.RemoteAddr = "127.0.0.1"

	service.SignUp_MockData.ExpectedError = nil

	props := map[string]interface{}{
		"email":      service.SignUp_MockData.Props.Email,
		"inviteCode": service.SignUp_MockData.Props.InviteCode,
		"language":   service.SignUp_MockData.Props.Language,
		"remoteAddr": service.SignUp_MockData.Props.RemoteAddr,
	}

	bytesRepresentation, err := json.Marshal(props)
	if err != nil {
		log.Fatalln(err)
	}

	res, err = http.Post(ts.URL+"/authorization/signup",
		"application/json;charset=utf-8",
		bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusAccepted {
		t.Errorf("Status code for /authorization/signup is wrong. Have: %d, want: %d.",
			res.StatusCode, http.StatusAccepted)
	}
}
