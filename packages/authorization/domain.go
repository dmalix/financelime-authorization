/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"github.com/dmalix/financelime-authorization/models"
	"net/http"
	"net/mail"
)

type APIMiddleware interface {
	RequestID(next http.Handler) http.Handler
	Authorization(next http.Handler) http.Handler
}

type Service interface {
	SignUp(email, language, inviteCode, remoteAddr string) error
	ConfirmUserEmail(confirmationKey string) (string, error)
	RequestAccessToken(email, password, clientID, remoteAddr string, device models.Device) (string, string, string, error)
	RefreshAccessToken(refreshToken, remoteAddr string) (string, string, string, error)
	GetListActiveSessions(encryptedUserData []byte) ([]models.Session, error)
}

type Repository interface {
	CreateUser(email, language, inviteCode, remoteAddr, confirmationKey string, inviteCodeRequired bool) error
	ConfirmUserEmail(confirmationKey string) (models.User, error)
	GetUserByAuth(email, password string) (models.User, error)
	GetUserByRefreshToken(RefreshToken string) (models.User, error)
	SaveSession(userID int64, publicSessionID, refreshToken, clientID, remoteAddr string, device models.Device) error
	UpdateSession(publicSessionID, refreshToken, remoteAddr string) error
	GetListActiveSessions(userID int64) ([]models.Session, error)
}

type Message interface {
	AddEmailMessageToQueue(messageQueue chan models.EmailMessage, to mail.Address, subject, body string, messageID ...string) error
}

type Cryptographer interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type Jwt interface {
	GenerateToken(publicSessionID string, userData []byte, tokenPurpose string, issuedAt ...int64) (string, error)
	VerifyToken(jwt string) (models.JwtData, error)
}
