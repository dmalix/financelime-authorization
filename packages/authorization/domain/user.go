/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package domain

import (
	"github.com/dmalix/financelime-rest-api/models"
	"net/mail"
)

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

type UserRepo interface {
	/*
		Create a new user
			----------------
			Return:
				confirmationID int64
				error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					PROPS_EMAIL:              param user.Email is not valid
					PROPS_INVITE:             parap user.InviteCode is not valid
					PROPS_LANG:               param user.Language is not valid
					PROPS_REMOTE_ADDR:        param remoteAddr is not valid
					PROPS_LINK_KEY:           param linkKey is not valid
					USER_ALREADY_EXIST:       a user with the email you specified already exists
					INVITE_NOT_EXIST_EXPIRED: the invite code does not exist or is expired
					INVITE_LIMIT:             the limit for issuing this invite code has been exhausted
	*/
	CreateUser(user *models.User, remoteAddr, linkKey string, inviteCodeRequired bool) (int64, error)
}

type UserSMTP interface {
	/*
		Send Email
			------------------
				Return:
					error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details])::
					        ------------------------------------------------
					        TLS_CONNECTION:  failed to Performing a TLS Connection
					        CONNECT_CLIENT:  failed to Connect New Client
					        CLIENT_AUTH:     failed to Perform client authentication
					        DATA_SENDER:     failed to Send the data of the sender
					        DATA_RECIPIENT:  failed to Send the data of the recipient
					        MESSAGE_DATA:    failed to Send message (DATA)
					        MESSAGE_WRITE:   failed to Send message (WRITE)
					        MESSAGE_CLOSE:   failed to Send message (CLOSE)
					        QUIT:            failed to Quit
	*/
	SendEmail(to mail.Address, subject, body string, messageID ...string) error
}
