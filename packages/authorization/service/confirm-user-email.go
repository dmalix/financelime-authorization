/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"net/mail"
	"strings"
)

/*
	   	Confirm user email
	   		----------------
	   		Return:
	   			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					CONFIRMATION_KEY_NOT_VALID: the confirmation key not valid
*/
// Related interfaces:
//	packages/authorization/domain.go
func (s *Service) ConfirmUserEmail(confirmationKey string) (string, error) {

	var (
		user     models.User
		err      error
		errLabel string
		message  string
	)

	user, err = s.repository.ConfirmUserEmail(confirmationKey)

	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case "PROPS_CONFIRMATION_KEY", "CONFIRMATION_KEY_NOT_FOUND_EXPIRED", "CONFIRMATION_KEY_ALREADY_CONFIRMED":
			return message,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					"CONFIRMATION_KEY_NOT_VALID",
					"the confirmation key not valid",
					err))

		default:
			errLabel = "gB7M0bJK"
			return message,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					errLabel,
					"a system error was returned",
					err))
		}
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: user.Email},
		s.languageContent.Data.User.Signup.Email.Password.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Signup.Email.Password.Body[s.languageContent.Language[user.Language]],
			user.Password),
		fmt.Sprintf(
			"<%s@%s>",
			user.Password,
			fmt.Sprintf("%s.%s", "confirm-user-email", s.config.DomainAPI)))

	if err != nil {
		errLabel = "XfCCWkb2"
		return message,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to send message to the user",
				err))
	}

	message = s.languageContent.Data.User.Signup.Page.Text[s.languageContent.Language[user.Language]]

	return message, nil
}
