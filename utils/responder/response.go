/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package responder

import (
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"log"
	"net/http"
	"strconv"
)

func Response(w http.ResponseWriter, r *http.Request, responseBody []byte, statusCode int, contentType ...string) {

	var (
		errorDetails error
		errorCode    int
	)
	message := "" /*fmt.Sprintf("%s %s %s",
	html.EscapeString(r.Method),
	html.EscapeString(r.URL.Path),
	html.EscapeString(r.Header.Get("request-id")))*/
	additionalInformation := fmt.Sprintf(
		"%s",
		strconv.Itoa(statusCode))

	if len(contentType) > 0 {
		w.Header().Set("content-type", contentType[0])
	}

	w.WriteHeader(statusCode)
	errorCode, errorDetails = w.Write(responseBody)
	if errorDetails != nil {
		log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
			fmt.Sprintf("Failed response (errorCode:%s): %s", strconv.Itoa(errorCode), message),
			errorDetails)
	}

	log.Printf("OK [%s: %s [%s]]", trace.GetCurrentPoint(), message, additionalInformation)

	return
}
