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
	"strings"
)

/*
	   	Refresh an access token
	   		----------------
	   		Return:
				publicSessionID string
				jwtAccess       string
				jwtRefresh      string
	   			err             error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
					------------------------------------------------
					INVALID_REFRESH_TOKEN:    Failed to validate the Refresh Token (JWT)
					USER_NOT_FOUND:           User is not found
*/
func (s *Service) RefreshAccessToken(refreshToken, remoteAddr string) (string, string, string, error) {

	var (
		err               error
		user              models.User
		sourceUserData    []byte
		encryptedUserData []byte
		jwtData           models.JwtData
		publicSessionID   string
		jwtAccess         string
		jwtRefresh        string
	)

	jwtData, err = s.jwt.VerifyToken(refreshToken)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"INVALID_REFRESH_TOKEN",
				"Failed to validate the Refresh Token (JWT)",
				err))
	}

	user, err = s.repository.GetUserByRefreshToken(refreshToken)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
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

	publicSessionID = jwtData.Payload.PublicSessionID

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

	err = s.repository.UpdateSession(publicSessionID, refreshToken, remoteAddr)
	if err != nil {
		return publicSessionID, jwtAccess, jwtRefresh,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to save the session",
				err))
	}

	return publicSessionID, jwtAccess, jwtRefresh, nil
}
