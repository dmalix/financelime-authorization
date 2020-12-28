/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/utils/jwt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"net/mail"
	"strings"
	"time"
)

/*
	   	Request an access token
	   		----------------
	   		Return:
				publicSessionID string
				jwtAccess       string
				jwtRefresh      string
	   			err             error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					PROPS:                    One or more of the input parameters are invalid
					USER_NOT_FOUND:           User is not found
*/
func (s *Service) RequestAccessToken(email, password,
	clientID, remoteAddr string, device models.Device) (string, string, string, error) {

	var (
		err               error
		user              models.User
		sourceUserData    []byte
		encryptedUserData []byte
		publicSessionID   string
		jwtAccess         string
		jwtRefresh        string
	)

	user, err = s.repository.GetUserByAuth(email, password)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case "PROPS_EMAIL", "PROPS_PASSWORD", "PROPS_LANG":
			return publicSessionID, jwtAccess, jwtRefresh,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					"PROPS",
					"One or more of the input parameters are invalid",
					err))
		case "USER_NOT_FOUND":
			return publicSessionID, jwtAccess, jwtRefresh,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					domainErrorCode,
					"User is not found",
					err))
		default:
			return publicSessionID, jwtAccess, jwtRefresh,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					trace.GetCurrentPoint(),
					"A system error was returned",
					err))
		}
	}

	publicSessionID, err = s.generatePublicID(user.ID)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate the publicSessionID value",
				err))
	}

	sourceUserData, err = json.Marshal(user)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to marshal the user struct",
				err))
	}

	encryptedUserData, err = s.cryptographer.Encrypt(sourceUserData)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to encrypt the user data",
				err))
	}

	jwtAccess, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.PropsPurposeAccess)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate an access token (JWT)",
				err))
	}

	jwtRefresh, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.PropsPurposeRefresh)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate an refresh token (JWT)",
				err))
	}

	err = s.repository.SaveSession(user.ID, publicSessionID, jwtRefresh, clientID, remoteAddr, device)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to save the session",
				err))
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: email},
		s.languageContent.Data.User.Login.Email.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Login.Email.Body[s.languageContent.Language[user.Language]],
			time.Now().UTC().String(),
			device.Platform,
			remoteAddr,
			s.config.DomainAPP),
		fmt.Sprintf(
			"<%s@%s>",
			remoteAddr,
			fmt.Sprintf("%s.%s", "request-access-token", s.config.DomainAPI)))
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to send message to the user",
				err))
	}

	return publicSessionID, jwtAccess, jwtRefresh, nil
}
