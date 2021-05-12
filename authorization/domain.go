/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"context"
	"net/http"
)

type API interface {
	signUp(ctx context.Context) http.Handler
	confirmUserEmail(ctx context.Context) http.Handler
	requestUserPasswordReset(ctx context.Context) http.Handler
	refreshAccessToken(ctx context.Context) http.Handler
	createAccessToken(ctx context.Context) http.Handler
	getListActiveSessions(ctx context.Context) http.Handler
	revokeRefreshToken(ctx context.Context) http.Handler
}

type Service interface {
	signUp(ctx context.Context, param serviceSignUpParam) error
	confirmUserEmail(ctx context.Context, confirmationKey string) (string, error)
	createAccessToken(ctx context.Context, param serviceCreateAccessTokenParam) (serviceAccessTokenReturn, error)
	refreshAccessToken(ctx context.Context, param serviceRefreshAccessTokenParam) (serviceAccessTokenReturn, error)
	revokeRefreshToken(ctx context.Context, param serviceRevokeRefreshTokenParam) error
	getListActiveSessions(ctx context.Context, encryptedUserData []byte) ([]session, error)
	requestUserPasswordReset(ctx context.Context, param serviceRequestUserPasswordResetParam) error
}

type Repository interface {
	createUser(ctx context.Context, param repoCreateUserParam) error
	confirmUserEmail(ctx context.Context, confirmationKey string) (user, error)
	getUserByAuth(ctx context.Context, param repoGetUserByAuthParam) (user, error)
	getUserByRefreshToken(ctx context.Context, refreshToken string) (user, error)
	saveSession(ctx context.Context, param repoSaveSessionParam) error
	updateSession(ctx context.Context, param repoUpdateSessionParam) error
	deleteSession(ctx context.Context, param repoDeleteSessionParam) error
	getListActiveSessions(ctx context.Context, userID int64) ([]session, error)
	requestUserPasswordReset(ctx context.Context, param repoRequestUserPasswordResetParam) (user, error)
}
