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
		account *models.Account
		linkKey string
		err     error
	)

	account = &models.Account{
		Email:      email,
		InviteCode: inviteCode,
		Language:   language,
	}

	linkKey = random.StringRand(16, 16, true)

	_, err = a.accountRepo.CreateAccount(account, remoteAddr, linkKey, a.inviteCodeRequired)
	if err != nil {
		customError := strings.Split(err.Error(), ":")[0]
		switch {
		case customError == "FL100":
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"lKJ1Qzfk",
				"Failed to create a new user",
				err.Error()))
		default:
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"4PtDRMCQ",
				"Failed to create a new user",
				err.Error()))
		}
	}

	return nil
}
