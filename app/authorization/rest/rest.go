/* Copyright © 2021. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package rest

import (
	"encoding/json"
	"github.com/dmalix/financelime-authorization/app/authorization"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"github.com/dmalix/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type rest struct {
	contextGetter middleware.ContextGetter
	service       authorization.Service
}

func NewREST(
	contextGetter middleware.ContextGetter,
	service authorization.Service) *rest {
	return &rest{
		contextGetter: contextGetter,
		service:       service,
	}
}

// SignUpStep1
// @Summary Create new user
// @Description The service sends a confirmation link to the specified email. After confirmation, the service will send a password for authorization.
// @ID signup_step1
// @Accept application/json;charset=utf-8
// @Param request-id header string true "RequestID"
// @Param model.SignUpRequest body model.SignUpRequest true "Data for creating a new user"
// @Success 204 "Successful operation"
// @Failure 400 {object} model.SignUpFailure400
// @Failure 404 {object} model.CommonFailure
// @Failure 409 {object} model.SignUpFailure409
// @Failure 500 {object} model.CommonFailure
// @Router /v1/user/ [post]
func (a *rest) SignUpStep1(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var requestInput model.SignUpRequest

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.DPanic("failed to read the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			logger.DPanic("failed to close a requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			logger.Error("failed to unmarshal the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, authorization.ErrorBadParams.Error(), http.StatusBadRequest)
			return
		}

		err = a.service.SignUpStep1(r.Context(), logger, model.ServiceSignUpParam{
			Email:      requestInput.Email,
			Language:   requestInput.Language,
			InviteCode: requestInput.InviteCode})
		if err != nil {
			logger.Error("failed to Sign Up", zap.Error(err), zap.String(requestIDKey, requestID))
			switch err {
			case authorization.ErrorBadParams:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			case authorization.ErrorUserAlreadyExist, authorization.ErrorInviteNotFound, authorization.ErrorInviteHasEnded:
				http.Error(w, err.Error(), http.StatusConflict)
				return
			default:
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// SignUpStep2
// @Summary Confirm User Email
// @Description API returns HTML-page with a message (success or error).
// @ID signup_step2
// @Produce text/plain;charset=utf-8
// @Param rid query string true "RequestID"
// @Param confirmationKey path string true "Confirmation Key"
// @Success 200 "Successful operation"
// @Failure 404 {object} model.CommonFailure
// @Failure 500 {object} model.CommonFailure
// @Router /u/{confirmationKey} [get]
func (a *rest) SignUpStep2(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		confirmationKey := vars["confirmationKey"]

		confirmationMessage, err := a.service.SignUpStep2(r.Context(), logger, confirmationKey)
		if err != nil {
			logger.Error("failed to confirm user email", zap.String(requestIDKey, requestID), zap.Error(err))
			switch err {
			case authorization.ErrorBadParamConfirmationKey, authorization.ErrorBadConfirmationKey:
				http.Error(w, statusMessageNotFound, http.StatusNotFound)
				return
			default:
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set(headerKeyContentType, headerValueTextPlain)
		w.WriteHeader(http.StatusOK)
		if code, err := w.Write([]byte(confirmationMessage)); err != nil {
			logger.DPanic("failed response", zap.Int("code", code), zap.Error(err), zap.String(requestIDKey, requestID))
			return
		}

		return
	})
}

// CreateAccessToken
// @Summary Create Access Token (Domain Action: Log In)
// @Description Create Access Token
// @ID create_access_token
// @Accept application/json;charset=utf-8
// @Produce application/json;charset=utf-8
// @Param request-id header string true "RequestID"
// @Param model.CreateAccessTokenRequest body model.CreateAccessTokenRequest true "Data for creating a new token"
// @Success 200 {object} model.AccessTokenResponse "Successful operation"
// @Failure 400 {object} model.CreateAccessTokenFailure400
// @Failure 404 {object} model.CreateAccessTokenFailure404
// @Failure 500 {object} model.CommonFailure
// @Router /v1/oauth/ [post]
func (a *rest) CreateAccessToken(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var requestInput model.CreateAccessTokenRequest
		var serviceAccessTokenReturn model.ServiceAccessTokenReturn

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.DPanic("failed to read the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		err = r.Body.Close()
		if err != nil {
			logger.DPanic("failed to close a requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			logger.Error("failed to unmarshal the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, authorization.ErrorBadParams.Error(), http.StatusBadRequest)
			return
		}

		serviceAccessTokenReturn, err =
			a.service.CreateAccessToken(r.Context(), logger, model.ServiceCreateAccessTokenParam{
				Email:     requestInput.Email,
				Password:  requestInput.Password,
				ClientID:  requestInput.ClientID,
				UserAgent: r.UserAgent(),
				Device:    requestInput.Device})

		if err != nil {
			logger.Error("failed to create an access token", zap.String(requestIDKey, requestID), zap.Error(err))
			switch err {
			case authorization.ErrorBadParams:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			case authorization.ErrorUserNotFound:
				http.Error(w, authorization.ErrorUserNotFound.Error(), http.StatusNotFound)
				return
			default:
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		responseBody, err := json.Marshal(model.AccessTokenResponse{
			PublicSessionID: serviceAccessTokenReturn.PublicSessionID,
			AccessJWT:       serviceAccessTokenReturn.AccessJWT,
			RefreshJWT:      serviceAccessTokenReturn.RefreshJWT,
		})
		if err != nil {
			logger.DPanic("failed to marshal model.AccessTokenResponse", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		w.Header().Set(headerKeyContentType, headerValueApplicationJson)
		w.WriteHeader(http.StatusOK)
		if code, err := w.Write(responseBody); err != nil {
			logger.DPanic("failed response", zap.Int("code", code), zap.Error(err),
				zap.String(requestIDKey, requestID))
			return
		}

		return
	})
}

// RefreshAccessToken
// @Summary Refresh Access Token (Domain Action: Renew authorization)
// @Description Refresh Access Token
// @ID refresh_access_token
// @Accept application/json;charset=utf-8
// @Produce application/json;charset=utf-8
// @Param request-id header string true "RequestID"
// @Param model.RefreshAccessTokenRequest body model.RefreshAccessTokenRequest true "Data for refreshing the access token"
// @Success 200 {object} model.AccessTokenResponse "Successful operation"
// @Failure 400 {object} model.RefreshAccessTokenFailure400
// @Failure 404 {object} model.RefreshAccessTokenFailure404
// @Failure 500 {object} model.CommonFailure
// @Router /v1/oauth/ [put]
func (a *rest) RefreshAccessToken(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var requestInput model.RefreshAccessTokenRequest
		var serviceAccessTokenReturn model.ServiceAccessTokenReturn

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.DPanic("failed to read the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			logger.DPanic("failed to close a requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			logger.Error("failed to unmarshal the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, authorization.ErrorBadParams.Error(), http.StatusBadRequest)
			return
		}

		serviceAccessTokenReturn, err = a.service.RefreshAccessToken(r.Context(), logger, requestInput.RefreshToken)
		if err != nil {
			logger.Error("failed to refresh an access token", zap.Error(err), zap.String(requestIDKey, requestID))
			switch err {
			case authorization.ErrorBadRefreshToken:
				http.Error(w, authorization.ErrorBadRefreshToken.Error(), http.StatusBadRequest)
			case authorization.ErrorUserNotFound:
				http.Error(w, authorization.ErrorUserNotFound.Error(), http.StatusNotFound)
				return
			default:
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		responseBody, err := json.Marshal(model.ServiceAccessTokenReturn{
			PublicSessionID: serviceAccessTokenReturn.PublicSessionID,
			AccessJWT:       serviceAccessTokenReturn.AccessJWT,
			RefreshJWT:      serviceAccessTokenReturn.RefreshJWT,
		})
		if err != nil {
			logger.DPanic("failed to marshal the access token", zap.Error(err), zap.String(requestIDKey, requestID))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", headerValueApplicationJson)
		w.WriteHeader(http.StatusOK)
		if code, err := w.Write(responseBody); err != nil {
			logger.DPanic("failed response", zap.Int("code", code), zap.Error(err),
				zap.String(requestIDKey, requestID))
			return
		}
		return
	})
}

// GetListActiveSessions
// @Summary Get a list of active sessions
// @Description Get a list of active sessions
// @ID get_list_active_sessions
// @Security authorization
// @Produce application/json;charset=utf-8
// @Param request-id header string true "RequestID"
// @Success 200 {object} []model.Session "Successful operation"
// @Failure 500 {object} model.CommonFailure
// @Router /v1/sessions/ [get]
func (a *rest) GetListActiveSessions(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var sessions []model.Session

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		accessTokenData, err := a.contextGetter.GetJwtData(r.Context())
		if err != nil {
			logger.DPanic("failed to get accessTokenData", zap.Error(err), zap.String(requestIDKey, requestID))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		sessions, err = a.service.GetListActiveSessions(r.Context(), logger, accessTokenData)
		if err != nil {
			logger.DPanic("failed to get the active sessions list", zap.Error(err), zap.String(requestIDKey, requestID))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		responseBody, err := json.Marshal(sessions)
		if err != nil {
			logger.DPanic("failed to marshal the active sessions list", zap.Error(err), zap.String(requestIDKey, requestID))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", headerValueApplicationJson)
		w.WriteHeader(http.StatusOK)
		if code, err := w.Write(responseBody); err != nil {
			logger.DPanic("failed response", zap.Int("code", code), zap.Error(err),
				zap.String(requestIDKey, requestID))
			return
		}
		return
	})
}

// RevokeRefreshToken
// @Summary Revoke Refresh Token (Domain Action: Log Out)
// @Description This request revoke the Refresh Token associated with the specified session. Thus, when the Access Token expires, then it cannot be renewed. And only after that, the user will be log out. Be aware that this query is idempotent.
// @ID revoke_refresh_token
// @Security authorization
// @Accept application/json;charset=utf-8
// @Param request-id header string true "RequestID"
// @Param model.RevokeRefreshTokenRequest body model.RevokeRefreshTokenRequest true "Data for revoking the Refresh Token"
// @Success 204 "Successful operation"
// @Failure 400 {object} model.RevokeRefreshTokenFailure400
// @Failure 500 {object} model.CommonFailure
// @Router /v1/session/ [delete]
func (a *rest) RevokeRefreshToken(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var requestInput model.RevokeRefreshTokenRequest
		var publicSessionID string

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.DPanic("failed to read the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}
		err = r.Body.Close()
		if err != nil {
			logger.DPanic("failed to close a requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			logger.Error("failed to unmarshal the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, authorization.ErrorBadParams.Error(), http.StatusBadRequest)
			return
		}

		accessTokenData, err := a.contextGetter.GetJwtData(r.Context())
		if err != nil {
			logger.DPanic("failed to get accessTokenData", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		if requestInput.PublicSessionID == "" {
			publicSessionID, err = a.contextGetter.GetJwtID(r.Context())
			if err != nil {
				logger.DPanic("failed to get publicSessionID", zap.Error(err))
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		} else {
			publicSessionID = requestInput.PublicSessionID
		}

		err = a.service.RevokeRefreshToken(r.Context(), logger, model.ServiceRevokeRefreshTokenParam{
			AccessTokenData: accessTokenData,
			PublicSessionID: publicSessionID,
		})
		if err != nil {
			logger.DPanic("failed to revoke the refresh token", zap.Error(err), zap.String(requestIDKey, requestID))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// ResetUserPasswordStep1
// @Summary Request to user password reset
// @Description The service sends a confirmation link to the specified email. After confirmation, the service will send a new password for authorization.
// @ID reset_user_password_step1
// @Accept application/json;charset=utf-8
// @Param request-id header string true "RequestID"
// @Param model.ResetUserPasswordRequest body model.ResetUserPasswordRequest true "Data for resetting your password"
// @Success 204 "Successful operation"
// @Failure 400 {object} model.RequestUserPasswordResetFailure400
// @Failure 404 {object} model.RequestUserPasswordResetFailure404
// @Failure 500 {object} model.CommonFailure
// @Router /v1/user/ [put]
func (a *rest) ResetUserPasswordStep1(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var requestInput model.ResetUserPasswordRequest

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.DPanic("failed to read the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		err = r.Body.Close()
		if err != nil {
			logger.DPanic("failed to close a requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			logger.Error("failed to unmarshal the requestBody", zap.String(requestIDKey, requestID), zap.Error(err))
			http.Error(w, authorization.ErrorBadParams.Error(), http.StatusBadRequest)
			return
		}

		err = a.service.ResetUserPasswordStep1(r.Context(), logger, requestInput.Email)
		if err != nil {
			logger.Error("failed to request a reset of the user's password", zap.Error(err), zap.String(requestIDKey, requestID))
			switch err {
			case authorization.ErrorBadParams:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			case authorization.ErrorUserNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			default:
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// ResetUserPasswordStep2
// @Summary Confirm to user password reset
// @Description API returns HTML-page with a message (success or error).
// @ID reset_user_password_step2
// @Produce text/plain;charset=utf-8
// @Param rid query string true "RequestID"
// @Param confirmationKey path string true "Confirmation Key"
// @Success 200 "Successful operation"
// @Failure 404 {object} model.CommonFailure
// @Failure 500 {object} model.CommonFailure
// @Router /p/{confirmationKey} [get]
func (a *rest) ResetUserPasswordStep2(logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		requestID, requestIDKey, err := a.contextGetter.GetRequestID(r.Context())
		if err != nil {
			logger.DPanic("failed to get requestID", zap.Error(err))
			http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
			return
		}

		confirmationKey := vars["confirmationKey"]

		confirmationMessage, err := a.service.ResetUserPasswordStep2(r.Context(), logger, confirmationKey)
		if err != nil {
			logger.Error("failed to confirm user password reset", zap.String(requestIDKey, requestID), zap.Error(err))
			switch err {
			case authorization.ErrorBadParamConfirmationKey, authorization.ErrorBadConfirmationKey:
				http.Error(w, statusMessageNotFound, http.StatusNotFound)
				return
			default:
				http.Error(w, statusMessageInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set(headerKeyContentType, headerValueTextPlain)
		w.WriteHeader(http.StatusOK)
		if code, err := w.Write([]byte(confirmationMessage)); err != nil {
			logger.DPanic("failed response", zap.Int("code", code), zap.Error(err), zap.String(requestIDKey, requestID))
			return
		}

		return
	})
}
