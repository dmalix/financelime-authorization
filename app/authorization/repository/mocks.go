/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"context"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"go.uber.org/zap"
)

type Mock struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

func (repo *Mock) CreateUser(_ context.Context, _ *zap.Logger, _ model.RepoCreateUserParam) error {
	return repo.Expected.Error
}

func (repo *Mock) ConfirmUserEmail(_ context.Context, _ *zap.Logger, _ string) (model.User, error) {
	var user model.User
	return user, repo.Expected.Error
}

func (repo *Mock) GetUserByAuth(_ context.Context, _ *zap.Logger, _ model.RepoGetUserByAuthParam) (model.User, error) {
	var user model.User
	return user, repo.Expected.Error
}

func (repo *Mock) SaveSession(_ context.Context, _ *zap.Logger, _ model.RepoSaveSessionParam) error {
	return repo.Expected.Error
}

func (repo *Mock) GetListActiveSessions(_ context.Context, _ *zap.Logger, _ int64) ([]model.Session, error) {
	var sessions []model.Session
	return sessions, repo.Expected.Error
}

func (repo *Mock) UpdateSession(_ context.Context, _ *zap.Logger, _ model.RepoUpdateSessionParam) error {
	return repo.Expected.Error
}

func (repo *Mock) GetUserByRefreshToken(_ context.Context, _ *zap.Logger, _ string) (model.User, error) {
	var user model.User
	return user, repo.Expected.Error
}

func (repo *Mock) DeleteSession(_ context.Context, _ *zap.Logger, _ model.RepoDeleteSessionParam) error {
	return repo.Expected.Error
}

func (repo *Mock) RequestUserPasswordReset(_ context.Context, _ *zap.Logger, _ model.RepoRequestUserPasswordResetParam) (model.User, error) {
	var user model.User
	return user, repo.Expected.Error
}
