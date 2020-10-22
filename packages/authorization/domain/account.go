/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package domain

import "github.com/dmalix/financelime-rest-api/models"

type AccountService interface {
	SignUp(email, inviteCode, language, remoteAddr string) error
}

type AccountRepo interface {
	CreateAccount(account *models.Account, remoteAddr, linkKey string, inviteCodeRequired bool) (int64, error)
}
