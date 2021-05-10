package authorization

import (
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

func (s *ServiceMockDescription) signUp(param serviceSignUpParam) error {

	if ServiceMockData.Expected.Error != nil {
		return ServiceMockData.Expected.Error
	}

	if param.email != ServiceMockData.Props.Email || param.inviteCode != ServiceMockData.Props.InviteCode ||
		param.language != ServiceMockData.Props.Language || param.remoteAddr != ServiceMockData.Props.RemoteAddr {
		return errors.New("DefaultError")
	}

	return ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) confirmUserEmail(_ string) (string, error) {
	return "", ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) createAccessToken(_ serviceCreateAccessTokenParam) (serviceAccessTokenReturn, error) {
	return serviceAccessTokenReturn{
		publicSessionID: "sessionID",
		accessJWT:       "accessToken",
		refreshJWT:      "refreshToken",
	}, ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) getListActiveSessions(encryptedUserData []byte) ([]session, error) {
	var sessions []session
	return sessions, ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) refreshAccessToken(_ serviceRefreshAccessTokenParam) (serviceAccessTokenReturn, error) {
	return serviceAccessTokenReturn{
			publicSessionID: "sessionID",
			accessJWT:       "accessToken",
			refreshJWT:      "refreshToken"},
		ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) revokeRefreshToken(_ serviceRevokeRefreshTokenParam) error {
	return ServiceMockData.Expected.Error
}

func (s *ServiceMockDescription) requestUserPasswordReset(_ serviceRequestUserPasswordResetParam) error {
	return ServiceMockData.Expected.Error
}
