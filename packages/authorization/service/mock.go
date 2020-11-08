/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
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
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (s *MockDescription) SignUp(email, inviteCode, language, remoteAddr string) error {

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
