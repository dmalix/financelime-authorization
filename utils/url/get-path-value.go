/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package url

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"html"
	"strings"
)

/*
	   	Get the value of a path element
	   		----------------
	   		Return:
				urlValue string
	   			error
*/
func GetPathValue(url string, valueNumber int) (string, error) {

	var (
		urlPathRequest string
		urlValues      []string
		urlValue       string
	)

	urlPathRequest = html.EscapeString(url)
	urlPathRequest = strings.ToLower(urlPathRequest)
	urlValues = strings.SplitN(urlPathRequest, "/", valueNumber+3)

	if len(urlValues) < (valueNumber + 2) {
		return urlValue,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(), "the value doesn't found", html.EscapeString(url)))
	}
	urlValue = urlValues[1+valueNumber]

	return urlValue, nil
}
