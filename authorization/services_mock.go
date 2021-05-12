package authorization

import (
	"context"
	"errors"
)

type incomingProps struct {
	Email      string
	InviteCode string
	Language   string
}

type ServiceMockDescription struct {
	Props struct {
		Email      string
		InviteCode string
		Language   string
		RemoteAddr string
	}
	Expected struct {
		Error error
	}
	repository RepoMockDescription
}

var ServiceMockData ServiceMockDescription

func (s *ServiceMockDescription) signUp(_ context.Context, param serviceSignUpParam) error {

	if ServiceMockData.Expected.Error != nil {
		return ServiceMockData.Expected.Error
	}

	if param.email != ServiceMockData.Props.Email || param.inviteCode != ServiceMockData.Props.InviteCode ||
		param.language != ServiceMockData.Props.Language || param.remoteAddr != ServiceMockData.Props.RemoteAddr {
		return errors.New("DefaultError")
	}

	return ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) confirmUserEmail(_ context.Context, _ string) (string, error) {
	return "", ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) createAccessToken(_ context.Context, _ serviceCreateAccessTokenParam) (serviceAccessTokenReturn, error) {
	return serviceAccessTokenReturn{
		publicSessionID: "sessionID",
		accessJWT:       "accessToken",
		refreshJWT:      "refreshToken",
	}, ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) getListActiveSessions(_ context.Context, _ []byte) ([]session, error) {
	var sessions []session
	return sessions, ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) refreshAccessToken(_ context.Context, _ serviceRefreshAccessTokenParam) (serviceAccessTokenReturn, error) {
	return serviceAccessTokenReturn{
			publicSessionID: "sessionID",
			accessJWT:       "accessToken",
			refreshJWT:      "refreshToken"},
		ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) revokeRefreshToken(_ context.Context, _ serviceRevokeRefreshTokenParam) error {
	return ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) requestUserPasswordReset(_ context.Context, _ serviceRequestUserPasswordResetParam) error {
	return ServiceMockData.Expected.Error
}
