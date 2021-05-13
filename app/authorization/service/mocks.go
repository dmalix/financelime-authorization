package service

import (
	"context"
	"errors"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	repository2 "github.com/dmalix/financelime-authorization/app/authorization/repository"
	"go.uber.org/zap"
)

type incomingProps struct {
	Email      string
	InviteCode string
	Language   string
}

type Mock struct {
	Props struct {
		Email      string
		InviteCode string
		Language   string
		RemoteAddr string
	}
	Expected struct {
		Error error
	}
	repository repository2.Mock
}

func (s *Mock) SignUp(_ context.Context, _ *zap.Logger, param model.ServiceSignUpParam) error {

	if s.Expected.Error != nil {
		return s.Expected.Error
	}

	if param.Email != s.Props.Email || param.InviteCode != s.Props.InviteCode ||
		param.Language != s.Props.Language {
		return errors.New("DefaultError")
	}

	return s.Expected.Error
}

func (s *Mock) ConfirmUserEmail(_ context.Context, _ *zap.Logger, _ string) (string, error) {
	return "", s.Expected.Error
}

func (s *Mock) CreateAccessToken(_ context.Context, _ *zap.Logger, _ model.ServiceCreateAccessTokenParam) (model.ServiceAccessTokenReturn, error) {
	return model.ServiceAccessTokenReturn{
		PublicSessionID: "sessionID",
		AccessJWT:       "accessToken",
		RefreshJWT:      "refreshToken",
	}, s.Expected.Error
}

func (s *Mock) GetListActiveSessions(_ context.Context, _ *zap.Logger, _ []byte) ([]model.Session, error) {
	var sessions []model.Session
	return sessions, s.Expected.Error
}

func (s *Mock) RefreshAccessToken(_ context.Context, _ *zap.Logger, _ string) (model.ServiceAccessTokenReturn, error) {
	return model.ServiceAccessTokenReturn{
			PublicSessionID: "sessionID",
			AccessJWT:       "accessToken",
			RefreshJWT:      "refreshToken"},
		s.Expected.Error
}

func (s *Mock) RevokeRefreshToken(_ context.Context, _ *zap.Logger, _ model.ServiceRevokeRefreshTokenParam) error {
	return s.Expected.Error
}

func (s *Mock) RequestUserPasswordReset(_ context.Context, _ *zap.Logger, _ string) error {
	return s.Expected.Error
}
