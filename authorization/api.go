/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func NewAPI(service Service) *api {
	return &api{
		service: service,
	}
}

const (
	contentTypeApplicationJson = "application/json;charset=utf-8"
	contentTypeTextPlain       = "text/plain;charset=utf-8"
)

// signUp
// @Summary Create new user
// @Description The service sends a confirmation link to the specified email. After confirmation, the service will send a password for authorization.
// @ID create_new_user
// @Accept application/json;charset=utf-8
// @Param apiSignUpRequest body apiSignUpRequest true "Data for creating a new user"
// @Success 204 "Successful operation"
// @Failure 400 {object} apiCommonFailure
// @Failure 404 {object} apiCommonFailure
// @Failure 409 {object} apiSignUpFailure409
// @Failure 500 {object} apiCommonFailure
// @Router /v1/signup [post]
func (api *api) signUp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			requestInput    apiSignUpRequest
			requestBody     []byte
			remoteAddr      string
			err             error
			domainErrorCode string
			errorMessage    string
		)

		if strings.ToLower(r.Header.Get("content-type")) != contentTypeApplicationJson {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Header 'content-type:%s' not found",
					contentTypeApplicationJson))
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		requestBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				fmt.Sprintf("Failed to get a requestBody"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}
		err = r.Body.Close()
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				fmt.Sprintf("Failed to close a requestBody"),
				err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				fmt.Sprintf("Failed to convert a requestBody requestInput to struct"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		remoteAddr = r.Header.Get("X-Real-IP")
		if len(remoteAddr) == 0 {
			remoteAddr = r.RemoteAddr
		}

		err = api.service.signUp(serviceSignUpParam{
			email:      requestInput.Email,
			language:   requestInput.Language,
			inviteCode: requestInput.InviteCode,
			remoteAddr: remoteAddr,
		})
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to Sign Up"
			switch domainErrorCode {
			case domainErrorCodeBadParams:
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, domainErrorCode, http.StatusBadRequest)
				return
			case domainErrorCodeUserAlreadyExist, domainErrorCodeInviteNotFound, domainErrorCodeInviteHasEnded:
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, domainErrorCode, http.StatusConflict)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, domainErrorCode, http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// confirmUserEmail
// @Summary Confirm User Email
// @Description API returns HTML-page with a message (success or error).
// @ID confirm_user_email
// @Produce text/plain;charset=utf-8
// @Param confirmationKey path string true "Confirmation Key"
// @Success 200 "Successful operation"
// @Failure 404 {object} apiCommonFailure
// @Failure 500 {object} apiCommonFailure
// @Router /v1/u/{confirmationKey} [get]
func (api *api) confirmUserEmail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			vars                = mux.Vars(r)
			confirmationKey     string
			confirmationMessage string
			err                 error
			domainErrorCode     string
			errorMessage        string
		)

		confirmationKey = vars["confirmationKey"]

		confirmationMessage, err = api.service.confirmUserEmail(confirmationKey)
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to confirm user email"
			switch domainErrorCode {
			case domainErrorCodeBadConfirmationKey: // the confirmation key not valid
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("content-type", contentTypeTextPlain)
		w.WriteHeader(http.StatusOK)
		if errorCode, err := w.Write([]byte(confirmationMessage)); err != nil {
			log.Printf("ERROR %s %s [%s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed response [errorCode:%s]", strconv.Itoa(errorCode)),
				err)
		}

		return
	})
}

// createAccessToken
// @Summary Create Access Token (Domain Action: Log In)
// @Description Create Access Token
// @ID create_access_token
// @Accept application/json;charset=utf-8
// @Produce application/json;charset=utf-8
// @Param apiCreateAccessTokenRequest body apiCreateAccessTokenRequest true "Data for creating a new token"
// @Success 200 {object} apiAccessTokenResponse "Successful operation"
// @Failure 400 {object} apiCommonFailure
// @Failure 404 {object} apiCommonFailure
// @Failure 500 {object} apiCommonFailure
// @Router /v1/oauth/token [post]
func (api *api) createAccessToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			requestInput      apiCreateAccessTokenRequest
			requestBody       []byte
			responseBody      []byte
			accessTokenReturn serviceAccessTokenReturn
			remoteAddr        string
			err               error
			domainErrorCode   string
			errorMessage      string
		)

		if strings.ToLower(r.Header.Get("content-type")) != contentTypeApplicationJson {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Header 'content-type:%s' not found",
					contentTypeApplicationJson))
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		requestBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get a body"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}
		err = r.Body.Close()
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to close a body"),
				err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to convert a body requestInput to struct"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		remoteAddr = r.Header.Get("X-Real-IP")
		if len(remoteAddr) == 0 {
			remoteAddr = r.RemoteAddr
		}

		accessTokenReturn, err =
			api.service.createAccessToken(serviceCreateAccessTokenParam{
				email:      requestInput.Email,
				password:   requestInput.Password,
				clientID:   requestInput.ClientID,
				remoteAddr: remoteAddr,
				userAgent:  r.UserAgent(),
				device:     requestInput.Device,
			})

		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to create an access token"
			switch domainErrorCode {
			case domainErrorCodeBadParams: // One or more of the input parameters are invalid
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, domainErrorCodeBadParams, http.StatusBadRequest)
				return
			case domainErrorCodeUserNotFound: // User is not found
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				w.Header().Add("domain-error-code", domainErrorCode)
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
				return
			}
		}

		responseBody, err = json.Marshal(apiAccessTokenResponse{
			PublicSessionID: accessTokenReturn.publicSessionID,
			AccessJWT:       accessTokenReturn.accessJWT,
			RefreshJWT:      accessTokenReturn.refreshJWT,
		})
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, "500 Server Internal Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", contentTypeApplicationJson)
		w.WriteHeader(http.StatusOK)
		if errorCode, err := w.Write(responseBody); err != nil {
			log.Printf("ERROR %s %s [%s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed response [errorCode:%s]", strconv.Itoa(errorCode)),
				err)
		}

		return
	})
}

// refreshAccessToken
// @Summary Refresh Access Token (Domain Action: Renew authorization)
// @Description Refresh Access Token
// @ID refresh_access_token
// @Accept application/json;charset=utf-8
// @Produce application/json;charset=utf-8
// @Param apiRefreshAccessTokenRequest body apiRefreshAccessTokenRequest true "Data for refreshing the access token"
// @Success 200 {object} apiAccessTokenResponse "Successful operation"
// @Failure 400 {object} apiCommonFailure
// @Failure 404 {object} apiCommonFailure
// @Failure 500 {object} apiCommonFailure
// @Router /v1/oauth/token [put]
func (api *api) refreshAccessToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			requestInput      apiRefreshAccessTokenRequest
			requestBody       []byte
			responseBody      []byte
			accessTokenReturn serviceAccessTokenReturn
			remoteAddr        string
			err               error
			domainErrorCode   string
			errorMessage      string
		)

		if strings.ToLower(r.Header.Get("content-type")) != contentTypeApplicationJson {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Header 'content-type:%s' not found",
					contentTypeApplicationJson))
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		requestBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				fmt.Sprintf("Failed to get a body"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}
		err = r.Body.Close()
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				fmt.Sprintf("Failed to close a body"),
				err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(requestBody, &requestInput)
		if err != nil {
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				fmt.Sprintf("Failed to convert a body requestInput to struct"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		remoteAddr = r.Header.Get("X-Real-IP")
		if len(remoteAddr) == 0 {
			remoteAddr = r.RemoteAddr
		}

		//response.PublicSessionID, response.AccessJWT, response.RefreshJWT, err =
		accessTokenReturn, err =
			api.service.refreshAccessToken(serviceRefreshAccessTokenParam{
				refreshToken: requestInput.RefreshToken,
				remoteAddr:   remoteAddr,
			})
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "failed to request an access token"
			switch domainErrorCode {
			case domainErrorCodeBadRefreshToken, domainErrorCodeUserNotFound: // One or more of the input parameters are invalid
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		responseBody, err = json.Marshal(apiAccessTokenResponse{
			PublicSessionID: accessTokenReturn.publicSessionID,
			AccessJWT:       accessTokenReturn.accessJWT,
			RefreshJWT:      accessTokenReturn.refreshJWT,
		})
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", contentTypeApplicationJson)
		w.WriteHeader(http.StatusOK)
		if errorCode, err := w.Write(responseBody); err != nil {
			log.Printf("ERROR %s %s [%s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed response [errorCode:%s]", strconv.Itoa(errorCode)),
				err)
		}
		return
	})
}

// revokeRefreshToken
// @Summary Revoke Refresh Token (Domain Action: Log Out)
// @Description This request revoke the Refresh Token associated with the specified session. Thus, when the Access Token expires, then it cannot be renewed. And only after that, the user will be log out. Be aware that this query is idempotent.
// @ID revoke_refresh_token
// @Security authorization
// @Accept application/json;charset=utf-8
// @Param apiRevokeRefreshTokenRequest body apiRevokeRefreshTokenRequest true "Data for revoking the Refresh Token"
// @Success 204 "Successful operation"
// @Failure 400 {object} apiCommonFailure
// @Failure 401 {object} apiCommonFailure
// @Failure 404 {object} apiCommonFailure
// @Failure 500 {object} apiCommonFailure
// @Router /v1/oauth/sessions [delete]
func (api *api) revokeRefreshToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			requestInput      apiRevokeRefreshTokenRequest
			err               error
			errorMessage      string
			body              []byte
			publicSessionID   string
			encryptedUserData []byte
		)

		if strings.ToLower(r.Header.Get("content-type")) != contentTypeApplicationJson {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Header 'content-type:%s' not found",
					contentTypeApplicationJson))
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get a body"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}
		err = r.Body.Close()
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to close a body"),
				err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &requestInput)
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to convert a body props to struct"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		if r.Context().Value(middleware.ContextEncryptedUserData) == nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get ContextEncryptedUserData from the request"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		encryptedUserData = r.Context().Value(middleware.ContextEncryptedUserData).([]byte)

		if len(requestInput.PublicSessionID) == 0 {
			if r.Context().Value(middleware.ContextPublicSessionID) == nil {
				// TODO Add RequestID from context
				// TODO Add PublicSessionID from context
				log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
					fmt.Sprintf("Failed to get ContextPublicSessionID from the request"),
					err)
				http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
				return
			}
			publicSessionID = r.Context().Value(middleware.ContextPublicSessionID).(string)
		} else {
			publicSessionID = requestInput.PublicSessionID
		}

		err = api.service.revokeRefreshToken(serviceRevokeRefreshTokenParam{
			encryptedUserData: encryptedUserData,
			publicSessionID:   publicSessionID,
		})
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// requestUserPasswordReset
// @Summary Request to user password reset
// @Description The service sends a confirmation link to the specified email. After confirmation, the service will send a new password for authorization.
// @ID user_password_reset
// @Accept application/json;charset=utf-8
// @Param apiRequestUserPasswordResetRequest body apiRequestUserPasswordResetRequest true "Data for resetting your password"
// @Success 204 "Successful operation"
// @Failure 400 {object} apiCommonFailure
// @Failure 404 {object} apiCommonFailure
// @Failure 500 {object} apiCommonFailure
// @Router /v1/resetpassword [post]
func (api *api) requestUserPasswordReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			requestInput    apiRequestUserPasswordResetRequest
			body            []byte
			remoteAddr      string
			err             error
			domainErrorCode string
			errorMessage    string
		)

		if strings.ToLower(r.Header.Get("content-type")) != contentTypeApplicationJson {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Header 'content-type:%s' not found",
					contentTypeApplicationJson))
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get a body"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		err = r.Body.Close()
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to close a body"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		err = json.Unmarshal(body, &requestInput)
		if err != nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to convert a body props to struct"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		remoteAddr = r.Header.Get("X-Real-IP")
		if len(remoteAddr) == 0 {
			remoteAddr = r.RemoteAddr
		}

		err = api.service.requestUserPasswordReset(serviceRequestUserPasswordResetParam{
			email:      requestInput.Email,
			remoteAddr: remoteAddr})
		if err != nil {
			domainErrorCode = strings.Split(err.Error(), ":")[0]
			errorMessage = "Failed to request a user password reset."
			switch domainErrorCode {
			case domainErrorCodeBadParams: // one or more of the input parameters are invalid
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, domainErrorCodeBadParams, http.StatusBadRequest)
				return
			case domainErrorCodeUserNotFound: // a user with the email specified not found
				log.Printf("ERROR [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
				return
			default:
				log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
				http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
		return
	})
}

// getListActiveSessions
// @Summary Get a list of active sessions
// @Description Get a list of active sessions
// @ID get_list_active_sessions
// @Security authorization
// @Produce application/json;charset=utf-8
// @Success 200 {object} []session "Successful operation"
// @Failure 401 {object} apiCommonFailure
// @Failure 404 {object} apiCommonFailure
// @Failure 500 {object} apiCommonFailure
// @Router /v1/oauth/sessions [get]
func (api *api) getListActiveSessions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			responseBody      []byte
			err               error
			errorMessage      string
			encryptedUserData []byte
			sessions          []session
		)

		if r.Context().Value(middleware.ContextEncryptedUserData) == nil {
			// TODO Add RequestID from context
			// TODO Add PublicSessionID from context
			log.Printf("ERROR [%s: %s [%s]]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed to get Context from the request"),
				err)
			http.Error(w, httpErrorTextNotFound, http.StatusNotFound)
			return
		}

		encryptedUserData = r.Context().Value(middleware.ContextEncryptedUserData).([]byte)

		sessions, err = api.service.getListActiveSessions(encryptedUserData)
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}

		responseBody, err = json.Marshal(sessions)
		if err != nil {
			log.Printf("FATAL [%s:%s[%s]]", trace.GetCurrentPoint(), errorMessage, err)
			http.Error(w, httpErrorTextInternalServerError, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", contentTypeApplicationJson)
		w.WriteHeader(http.StatusOK)
		if errorCode, err := w.Write(responseBody); err != nil {
			log.Printf("ERROR %s %s [%s]", trace.GetCurrentPoint(),
				fmt.Sprintf("Failed response [errorCode:%s]", strconv.Itoa(errorCode)),
				err)
		}
		return
	})
}
