/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package router

import (
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
	exampleRouter(mux)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// StatusNotFound
	// --------------

	res, err = http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code for /not-exists is wrong. Have: %d, want: %d.",
			res.StatusCode, http.StatusNotFound)
	}

	// StatusMethodNotAllowed
	// ----------------------

	res, err = http.Get(ts.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status code for /test is wrong. Have: %d, want: %d.",
			res.StatusCode, http.StatusMethodNotAllowed)
	}

	// StatusAccepted
	// --------------

	res, err = http.Post(ts.URL+"/test",
		"application/json;charset=utf-8",
		nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for /test is wrong. Have: %d, want: %d.",
			res.StatusCode, http.StatusOK)
	}
}
