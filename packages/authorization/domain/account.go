/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package domain

import "github.com/dmalix/financelime-rest-api/models"

type AccountService interface {
	SignUp(email, inviteCode, language, remoteAddr string) error
}

type AccountRepo interface {
	/*
		Create a new account
			----------------
			Return:
				confirmationID int64
				error  - system or custom error (format FLNNN:[details]):
				         ------------------------------------------------
				         FL100 - Param account.Email is not valid
				         FL101 - Parap account.InviteCode is not valid
				         FL102 - Param account.Language is not valid
				         FL103 - An account with the specified email already exists
				         FL104 - The invite code does not exist or is expired
				         FL105 - The limit for issuing this invite code has been exhausted
				         FL106 - Param remoteAddr is not valid
	*/
	CreateAccount(account *models.Account, remoteAddr, linkKey string, inviteCodeRequired bool) (int64, error)
}
