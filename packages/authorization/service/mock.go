/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"fmt"
)

type AuthorizationServiceMock struct{
	Props struct {
		SignUp struct {
			Email      string
			InviteCode string
			Language   string
			RemoteAddr string
		}
	}
}

//noinspection GoNameStartsWithPackageName
var ServiceMockValue AuthorizationServiceMock

func (a *AuthorizationServiceMock) SignUp(email, inviteCode, language, remoteAddr string) error {

	const (
		theEmailExists = "email.exists@financelime.com"
		theInviteCodeDoesNotExistOrIsExpired = "InviteCodeErrorFL104"
		theLimitForIssuingThisInviteCodeHasBeenExhausted = "InviteCodeErrorFL105"
		theParamRemoteAddrIsNotValid = "ParamRemoteAddrIsNotValid"
	)

	if email == theEmailExists {
		return errors.New(fmt.Sprintf("FL%s:[account.Email=%s]", "103", email))
	}
	if inviteCode == theInviteCodeDoesNotExistOrIsExpired {
		return errors.New(fmt.Sprintf("FL%s:[account.inviteCode=%s]", "104", inviteCode))
	}
	if inviteCode == theLimitForIssuingThisInviteCodeHasBeenExhausted {
		return errors.New(fmt.Sprintf("FL%s:[account.inviteCode=%s]", "105", inviteCode))
	}
	if remoteAddr == theParamRemoteAddrIsNotValid {
		return errors.New(fmt.Sprintf("FL%s:[account.inviteCode=%s]", "106", inviteCode))
	}

	if email != ServiceMockValue.Props.SignUp.Email {
		return errors.New(fmt.Sprintf("FL%s:[account.Email=%s]", "100", email))
	}
	if inviteCode != ServiceMockValue.Props.SignUp.InviteCode {
		return errors.New(fmt.Sprintf("FL%s:[account.InviteCode=%s]", "101", inviteCode))
	}
	if language != ServiceMockValue.Props.SignUp.Language {
		return errors.New(fmt.Sprintf("FL%s:[account.Language=%s]", "102", inviteCode))
	}


	return nil
}
