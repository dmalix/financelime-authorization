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
	"net/mail"
	"strings"
	"time"
)

/*
	   	Request an access token
	   		----------------
	   		Return:
				jwtAccess   string
				jwtRefresh  string
	   			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					PROPS:                    One or more of the input parameters are invalid
					USER_NOT_FOUND:           User is not found
*/
func (s *Service) RequestAccessToken(email, password,
	clientID, remoteAddr string, device models.Device) (string, string, error) {

	var (
		err      error
		errLabel string

		user              models.User
		publicSessionID   string
		sourceUserData    []byte
		encryptedUserData []byte

		accessToken  string
		refreshToken string
	)

	user, err = s.repository.GetUserByAuth(email, password)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case "PROPS_EMAIL", "PROPS_PASSWORD", "PROPS_LANG":
			return accessToken, refreshToken,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					"PROPS",
					"One or more of the input parameters are invalid",
					err))
		case "USER_NOT_FOUND":
			return accessToken, refreshToken,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					domainErrorCode,
					"User is not found",
					err))
		default:
			errLabel = "3REuo5jS"
			return accessToken, refreshToken,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					errLabel,
					"A system error was returned",
					err))
		}
	}

	publicSessionID, err = s.generatePublicID(user.ID)
	if err != nil {
		errLabel = "o53REujS"
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to generate the publicSessionID value",
				err))
	}

	sourceUserData, err = json.Marshal(user)
	if err != nil {
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to marshal the user struct",
				err))
	}

	encryptedUserData, err = s.cryptographer.Encrypt(sourceUserData)
	if err != nil {
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to marshal the user struct",
				err))
	}

	accessToken, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.PropsPurposeAccess)
	if err != nil {
		errLabel = "Ro53EujS"
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to generate an access token (JWT)",
				err))
	}

	refreshToken, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.PropsPurposeRefresh)
	if err != nil {
		errLabel = "D8JVbpWO"
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to generate an refresh token (JWT)",
				err))
	}

	err = s.repository.SaveSession(user.ID, publicSessionID, clientID, remoteAddr, device)
	if err != nil {
		errLabel = "6QqPfJGg"
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
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
		errLabel = "XfCCWkb2"
		return accessToken, refreshToken,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to send message to the user",
				err))
	}

	return accessToken, refreshToken, nil
}
