/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package rest

import (
	"encoding/json"
	"github.com/dmalix/financelime-authorization/app/information"
	"github.com/dmalix/financelime-authorization/app/information/model"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"go.uber.org/zap"
	"net/http"
)

type rest struct {
	contextGetter middleware.ContextGetter
	Service       information.Service
}

func NewREST(
	contextGetter middleware.ContextGetter,
	service information.Service) *rest {
	return &rest{
		contextGetter: contextGetter,
		Service:       service,
	}
}

// Version
// @Summary Get the Service version
// @Description Get Version
// @id get_version
// @Param request-id header string true "RequestID"
// @Produce application/json;charset=utf-8
// @Success 200 {object} model.VersionResponse "Successful operation"
// @Router /v1/version [get]
func (rest *rest) Version(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var versionResponse model.VersionResponse

		requestID, requestIDKey, err := rest.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		versionResponse.Number, versionResponse.Build, err = rest.Service.Version(r.Context(), logger)
		if err != nil {
			logger.DPanic("failed to get version", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		versionResponseJSON, err := json.Marshal(&versionResponse)
		if err != nil {
			logger.DPanic("failed to unmarshal the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		w.Header().Set(headerKeyContentType, headerValueApplicationJson)
		w.WriteHeader(http.StatusOK)
		if code, err := w.Write(versionResponseJSON); err != nil {
			logger.DPanic("failed response", zap.Int("code", code), zap.Error(err),
				zap.String(requestIDKey, requestID))
			return
		}

		return
	})
}
