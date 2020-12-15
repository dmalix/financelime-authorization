/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/utils/random"
	"net/mail"
	"strings"
)

/*
	   	Request user password reset
	   		----------------
	   		Return:
	   			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					PROPS:                one or more of the input parameters are invalid
					USER_NOT_FOUND:       a user with the email specified not found
*/
func (s *Service) RequestUserPasswordReset(email, remoteAddr string) error {

	var (
		confirmationKey string
		err             error
		errLabel        string
		user            models.User
	)

	confirmationKey = random.StringRand(16, 16, true)

	user, err = s.repository.RequestUserPasswordReset(email, remoteAddr, confirmationKey)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case "PROPS_EMAIL", "PROPS_REMOTE_ADDR", "PROPS_CONFIRMATION_KEY":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS",
				"one or more of the input parameters are invalid",
				err))
		case "USER_NOT_FOUND": // a user with the email specified not found
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"a user with the email you specified already exists",
				err))
		default:
			errLabel = "ahXah3qu"
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"a system error was returned",
				err))
		}
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: email},
		s.languageContent.Data.User.ResetPassword.Email.Request.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.ResetPassword.Email.Request.Body[s.languageContent.Language[user.Language]],
			remoteAddr, s.config.DomainAPI, confirmationKey),
		fmt.Sprintf(
			"<%s@%s>",
			confirmationKey,
			fmt.Sprintf("%s.%s", "reset-password", s.config.DomainAPI)))
	if err != nil {
		errLabel = "bXfCCWk2"
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to send message to the user",
			err))
	}

	return nil
}
