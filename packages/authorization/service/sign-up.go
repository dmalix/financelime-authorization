/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/utils/random"
	"strings"
)

func (a *Service) SignUp(email, inviteCode, language, remoteAddr string) error {

	var (
		user     *models.User
		linkKey  string
		err      error
		errLabel string
	)

	user = &models.User{
		Email:      email,
		InviteCode: inviteCode,
		Language:   language,
	}

	linkKey = random.StringRand(16, 16, true)

	_, err = a.userRepo.CreateUser(user, remoteAddr, linkKey, a.inviteCodeRequired)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case "PROPS_EMAIL", "PROPS_INVITE_CODE", "PROPS_LANG", "PROPS_REMOTE_ADDR", "PROPS_LINK_KEY":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS",
				"one or more of the input parameters are invalid",
				err))
		case "USER_ALREADY_EXIST":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"a user with the email you specified already exists",
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

	return nil
}
