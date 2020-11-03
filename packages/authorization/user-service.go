/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

type UserService interface {
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
	SignUp(email, inviteCode, language, remoteAddr string) error
}
