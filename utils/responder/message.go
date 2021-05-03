/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package responder

import (
	"fmt"
	"html"
	"net/http"
)

func Message(request *http.Request) string {
	var message string
	if request != nil {
		message = fmt.Sprintf(
			"%s %s",
			html.EscapeString(request.Method),
			html.EscapeString(request.URL.Path))
	}
	return message
}
