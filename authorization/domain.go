/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"net/http"
)

type APIMiddleware interface {
	requestID(next http.Handler) http.Handler
	authorization(next http.Handler) http.Handler
}

type API interface {
	signUp() http.Handler
	confirmUserEmail() http.Handler
	requestUserPasswordReset() http.Handler
	refreshAccessToken() http.Handler
	createAccessToken() http.Handler
	getListActiveSessions() http.Handler
	revokeRefreshToken() http.Handler
}

type Service interface {
	signUp(serviceSignUpParam) error
	confirmUserEmail(confirmationKey string) (string, error)
	createAccessToken(serviceCreateAccessTokenParam) (serviceAccessTokenReturn, error)
	refreshAccessToken(serviceRefreshAccessTokenParam) (serviceAccessTokenReturn, error)
	revokeRefreshToken(serviceRevokeRefreshTokenParam) error
	getListActiveSessions(encryptedUserData []byte) ([]session, error)
	requestUserPasswordReset(serviceRequestUserPasswordResetParam) error
}

type Repository interface {
	createUser(repoCreateUserParam) error
	confirmUserEmail(confirmationKey string) (user, error)
	getUserByAuth(repoGetUserByAuthParam) (user, error)
	getUserByRefreshToken(refreshToken string) (user, error)
	saveSession(repoSaveSessionParam) error
	updateSession(repoUpdateSessionParam) error
	deleteSession(repoDeleteSessionParam) error
	getListActiveSessions(userID int64) ([]session, error)
	requestUserPasswordReset(repoRequestUserPasswordResetParam) (user, error)
}
