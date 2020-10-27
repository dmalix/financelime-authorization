/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
)

type MockType struct{
	Values struct {
		SignUp struct {
			Props struct {
				Email      string
				InviteCode string
				Language   string
				RemoteAddr string
			}
			ExpectedError error
		}
	}
}

var Mock MockType

func (a *MockType) SignUp(email, inviteCode, language, remoteAddr string) error {

	if Mock.Values.SignUp.ExpectedError != nil {
		return  Mock.Values.SignUp.ExpectedError
	}

	if email != Mock.Values.SignUp.Props.Email && inviteCode != Mock.Values.SignUp.Props.InviteCode &&
		language != Mock.Values.SignUp.Props.Language && remoteAddr != Mock.Values.SignUp.Props.RemoteAddr {
		return errors.New("DefaultError")
	}

	return nil
}
