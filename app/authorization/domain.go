/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"context"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"go.uber.org/zap"
	"net/http"
)

type REST interface {
	SignUp(logger *zap.Logger) http.Handler
	ConfirmUserEmail(logger *zap.Logger) http.Handler
	RequestUserPasswordReset(logger *zap.Logger) http.Handler
	CreateAccessToken(logger *zap.Logger) http.Handler
	RefreshAccessToken(logger *zap.Logger) http.Handler
	GetListActiveSessions(logger *zap.Logger) http.Handler
	RevokeRefreshToken(logger *zap.Logger) http.Handler
}

type Service interface {
	SignUp(ctx context.Context, logger *zap.Logger, param model.ServiceSignUpParam) error
	ConfirmUserEmail(ctx context.Context, logger *zap.Logger, confirmationKey string) (string, error)
	CreateAccessToken(ctx context.Context, logger *zap.Logger, param model.ServiceCreateAccessTokenParam) (model.ServiceAccessTokenReturn, error)
	RefreshAccessToken(ctx context.Context, logger *zap.Logger, refreshToken string) (model.ServiceAccessTokenReturn, error)
	RevokeRefreshToken(ctx context.Context, logger *zap.Logger, param model.ServiceRevokeRefreshTokenParam) error
	GetListActiveSessions(ctx context.Context, logger *zap.Logger, encryptedUserData []byte) ([]model.Session, error)
	RequestUserPasswordReset(ctx context.Context, logger *zap.Logger, email string) error
}

type Repository interface {
	CreateUser(ctx context.Context, logger *zap.Logger, param model.RepoCreateUserParam) error
	ConfirmUserEmail(ctx context.Context, logger *zap.Logger, confirmationKey string) (model.User, error)
	GetUserByAuth(ctx context.Context, logger *zap.Logger, param model.RepoGetUserByAuthParam) (model.User, error)
	GetUserByRefreshToken(ctx context.Context, logger *zap.Logger, refreshToken string) (model.User, error)
	SaveSession(ctx context.Context, logger *zap.Logger, param model.RepoSaveSessionParam) error
	UpdateSession(ctx context.Context, logger *zap.Logger, param model.RepoUpdateSessionParam) error
	DeleteSession(ctx context.Context, logger *zap.Logger, param model.RepoDeleteSessionParam) error
	GetListActiveSessions(ctx context.Context, logger *zap.Logger, userID int64) ([]model.Session, error)
	RequestUserPasswordReset(ctx context.Context, logger *zap.Logger, param model.RepoRequestUserPasswordResetParam) (model.User, error)
}
