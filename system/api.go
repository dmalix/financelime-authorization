/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package system

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"log"
	"net/http"
	"strconv"
)

type api struct {
	Service Service
}

func NewAPI(service Service) *api {
	return &api{
		Service: service,
	}
}

const (
	contentTypeApplicationJson = "application/json;charset=utf-8"
	contentTypeTextPlain       = "text/plain;charset=utf-8"
)

// version
// @Summary Get the Service version
// @Description Get Version
// @id get_version
// @Produce application/json;charset=utf-8
// @Success 200 {object} versionResponse "Successful operation"
// @Router /v1/version [get]
func (api *api) version() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			err                 error
			errorMessage        string
			versionResponse     versionResponse
			versionResponseJSON []byte
		)

		versionResponse.Number, versionResponse.Build, err = api.Service.version()
		if err != nil {
			errorMessage = "failed to get version"
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		versionResponseJSON, err = json.Marshal(&versionResponse)
		if err != nil {
			errorMessage = "failed to convert the version data to JSON-format"
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", contentTypeApplicationJson)
		w.WriteHeader(http.StatusOK)
		if errorCode, err := w.Write(versionResponseJSON); err != nil {
			log.Printf("ERROR %s %s [%s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed response [errorCode:%s]", strconv.Itoa(errorCode)),
				err)
		}

		return
	})
}
