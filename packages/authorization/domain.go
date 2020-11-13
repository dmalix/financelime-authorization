/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"github.com/dmalix/financelime-authorization/models"
	"net/mail"
)

type Service interface {
	SignUp(email, language, inviteCode, remoteAddr string) error
	ConfirmUserEmail(confirmationKey string) (string, error)
	RequestAccessToken(email, password, clientID, remoteAddr string, device models.Device) (string, string, error)
}

type Repository interface {
	CreateUser(email, language, inviteCode, remoteAddr, confirmationKey string, inviteCodeRequired bool) error
	ConfirmUserEmail(confirmationKey string) (models.User, error)
	GetUserByAuth(email, password string) (models.User, error)
	SaveSession(userID int64, publicSessionID, client_id, remoteAddr string, device models.Device) error
}

type Message interface {
	AddEmailMessageToQueue(messageQueue chan models.EmailMessage, to mail.Address, subject, body string, messageID ...string) error
}

type Jwt interface {
	GenerateToken(publicSessionID, tokenPurpose string, issuedAt ...int64) (string, error)
	VerifyToken(jwt string) (models.JwtData, error)
}
