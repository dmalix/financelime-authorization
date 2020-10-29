/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package domain

import "github.com/dmalix/financelime-rest-api/models"

type UserService interface {
	/*
	   	Create a new user
	   		----------------
	   		Return:
	   			error  - system or domain error code (format domain_error_code:description[details]):
					------------------------------------------------
					a1: one or more of the input parameters are invalid
					b1: a user with the email you specified already exists
					b2: the invite code does not exist or is expired
					b3: the limit for issuing this invite code has been exhausted
	*/
	SignUp(email, inviteCode, language, remoteAddr string) error
}

type UserRepo interface {
	/*
		Create a new user
			----------------
			Return:
				confirmationID int64
				error  - system or domain error code (format domain_error_code:description[details]):
					------------------------------------------------
					a1: param user.Email is not valid
					a2: parap user.InviteCode is not valid
					a3: param user.Language is not valid
					a4: param remoteAddr is not valid
					a5: param linkKey is not valid
					b1: a user with the email you specified already exists
					b2: the invite code does not exist or is expired
					b3: the limit for issuing this invite code has been exhausted
	*/
	CreateUser(user *models.User, remoteAddr, linkKey string, inviteCodeRequired bool) (int64, error)
}
