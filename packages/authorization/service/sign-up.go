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
		case "a1", "a2", "a3", "a4", "a5":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"a1",
				"one or more of the input parameters are invalid",
				err))
		case "b1":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"b1",
				"a user with the email you specified already exists",
				err))
		case "b2":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"b2",
				"the invite code does not exist or is expired",
				err))
		case "b3":
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"b3",
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
