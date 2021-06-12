package repository

import (
	"context"
	"github.com/dmalix/authorization-service/app/authorization/model"
	"go.uber.org/zap"
)

type Mock struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

func (repo *Mock) SignUpStep1(_ context.Context, _ *zap.Logger, _ model.RepoSignUpParam) error {
	return repo.Expected.Error
}

func (repo *Mock) SignUpStep2(_ context.Context, _ *zap.Logger, _ string) (model.User, error) {
	return model.User{}, repo.Expected.Error
}

func (repo *Mock) GetUserByAuth(_ context.Context, _ *zap.Logger, _ model.RepoGetUserByAuthParam) (model.User, error) {
	return model.User{}, repo.Expected.Error
}

func (repo *Mock) CreateSession(_ context.Context, _ *zap.Logger, _ model.RepoCreateSessionParam) error {
	return repo.Expected.Error
}

func (repo *Mock) GetListActiveSessions(_ context.Context, _ *zap.Logger, _ int64) ([]model.Session, error) {
	return nil, repo.Expected.Error
}

func (repo *Mock) UpdateSession(_ context.Context, _ *zap.Logger, _ model.RepoUpdateSessionParam) error {
	return repo.Expected.Error
}

func (repo *Mock) GetUserByRefreshToken(_ context.Context, _ *zap.Logger, _ string) (model.User, error) {
	return model.User{}, repo.Expected.Error
}

func (repo *Mock) DeleteSession(_ context.Context, _ *zap.Logger, _ model.RepoDeleteSessionParam) error {
	return repo.Expected.Error
}

func (repo *Mock) ResetUserPasswordStep1(_ context.Context, _ *zap.Logger, _ model.RepoResetUserPasswordParam) (model.User, error) {
	return model.User{}, repo.Expected.Error
}

func (repo *Mock) ResetUserPasswordStep2(_ context.Context, _ *zap.Logger, _ string) (model.User, error) {
	return model.User{}, repo.Expected.Error
}
