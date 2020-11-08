/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/utils/random"
	"net/mail"
	"strings"
)

/*
	   	Create a new user
	   		----------------
	   		Return:
	   			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					PROPS:                    one or more of the input parameters are invalid
					USER_ALREADY_EXIST:       a user with the email you specified already exists
					INVITE_NOT_EXIST_EXPIRED: the invite code does not exist or is expired
					INVITE_LIMIT:             the limit for issuing this invite code has been exhausted
*/
// Related interfaces:
//	packages/authorization/domain.go
func (s *Service) SignUp(propsEmail, propsInviteCode, propsLanguage, propsRemoteAddr string) error {

	var (
		user            *models.User
		confirmationKey string
		err             error
		errLabel        string
	)

	user = &models.User{
		Email:    propsEmail,
		Language: propsLanguage,
	}

	confirmationKey = random.StringRand(16, 16, true)

	err = s.repository.CreateUser(user, propsInviteCode, propsRemoteAddr, confirmationKey, s.inviteCodeRequired)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case "PROPS_EMAIL", "PROPS_INVITE_CODE", "PROPS_LANG", "PROPS_REMOTE_ADDR", "PROPS_CONFIRMATION_KEY":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS",
				"one or more of the input parameters are invalid",
				err))
		case "USER_ALREADY_EXIST":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"a user with the propsEmail you specified already exists",
				err))
		case "INVITE_NOT_EXIST_EXPIRED":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"the invite code does not exist or is expired",
				err))
		case "INVITE_LIMIT":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"the limit for issuing this invite code has been exhausted",
				err))

		default:
			errLabel = "4PtDRMCQ"
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"a system error was returned",
				err))
		}
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: propsEmail},
		s.languageContent.Data.User.Signup.Email.Confirm.Subject[s.languageContent.Language[propsLanguage]],
		fmt.Sprintf(
			s.languageContent.Data.User.Signup.Email.Confirm.Body[s.languageContent.Language[propsLanguage]],
			s.domainAPI, confirmationKey),
		fmt.Sprintf(
			"<%s@%s>",
			confirmationKey,
			fmt.Sprintf("%s.%s", "sign-up", s.domainAPI)))

	if err != nil {
		errLabel = "XfCCWkb2"
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to send message to the user",
			err))
	}

	return nil
}
