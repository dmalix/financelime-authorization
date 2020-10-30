/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
)

//noinspection GoSnakeCaseUsage
type SignUp_MockDescription struct {
	Props struct {
		Email      string
		InviteCode string
		Language   string
		RemoteAddr string
	}
	ExpectedError error
}

//noinspection GoSnakeCaseUsage
var SignUp_MockData SignUp_MockDescription

func (a *SignUp_MockDescription) SignUp(email, inviteCode, language, remoteAddr string) error {

	if SignUp_MockData.ExpectedError != nil {
		return SignUp_MockData.ExpectedError
	}

	if email != SignUp_MockData.Props.Email && inviteCode != SignUp_MockData.Props.InviteCode &&
		language != SignUp_MockData.Props.Language && remoteAddr != SignUp_MockData.Props.RemoteAddr {
		return errors.New("DefaultError")
	}

	return nil
}
