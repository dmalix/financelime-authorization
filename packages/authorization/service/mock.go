/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/packages/authorization/repository"
)

//noinspection GoSnakeCaseUsage
type MockDescription struct {
	Props struct {
		Email      string
		InviteCode string
		Language   string
		RemoteAddr string
	}
	Expected struct {
		Error error
	}
	repository repository.MockDescription
}

var MockData MockDescription

func (s *MockDescription) SignUp(email, language, inviteCode, remoteAddr string) error {

	if MockData.Expected.Error != nil {
		return MockData.Expected.Error
	}

	if email != MockData.Props.Email || inviteCode != MockData.Props.InviteCode ||
		language != MockData.Props.Language || remoteAddr != MockData.Props.RemoteAddr {
		return errors.New("DefaultError")
	}

	return MockData.Expected.Error
}

func (s *MockDescription) ConfirmUserEmail(confirmationKey string) (string, error) {
	return "", MockData.Expected.Error
}

func (s *MockDescription) RequestAccessToken(email, password, clientID, remoteAddr string, device models.Device) (string, string, string, error) {
	return "sessionID", "accessToken", "refreshToken", MockData.Expected.Error
}

func (s *MockDescription) GetListActiveSessions(encryptedUserData []byte) ([]models.Session, error) {
	var sessions []models.Session
	return sessions, MockData.Expected.Error
}

func (s *MockDescription) RefreshAccessToken(refreshToken, remoteAddr string) (string, string, string, error) {
	return "sessionID", "accessToken", "refreshToken", MockData.Expected.Error
}

func (s *MockDescription) RevokeRefreshToken(encryptedUserData []byte, publicSessionID string) error {
	return MockData.Expected.Error
}

func (s *MockDescription) RequestUserPasswordReset(email string, remoteAddr string) error {
	return MockData.Expected.Error
}
