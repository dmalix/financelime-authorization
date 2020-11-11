/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"strings"
	"time"
)

func (token *Token) VerifyToken(jwt string) (models.JwtData, error) {

	var (
		valueByte   []byte
		err         error
		errLabel    string
		jwtTokenArr []string
		jwtToken    string
		lifeTime    int
		jwtData     models.JwtData
	)

	const InvalidJwtToken = "Invalid JWT Token"
	const NoPadding rune = -1

	jwtTokenArr = strings.Split(jwt, ".")

	// Check Headers
	valueByte, err = base64.URLEncoding.WithPadding(NoPadding).DecodeString(jwtTokenArr[0])
	if err != nil {
		errLabel = "zp0Bj5Ao"
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", errLabel, InvalidJwtToken, err))
	}

	err = json.Unmarshal(valueByte, &jwtData.Headers)
	if err != nil {
		errLabel = "dXZbnxN0"
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", errLabel, InvalidJwtToken, err))
	}

	// Check Payload
	valueByte, err = base64.URLEncoding.WithPadding(NoPadding).DecodeString(jwtTokenArr[1])
	if err != nil {
		errLabel = "liXr3xoK"
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", errLabel, InvalidJwtToken, err))
	}
	err = json.Unmarshal(valueByte, &jwtData.Payload)
	if err != nil {
		errLabel = "4PqkJfm6"
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", errLabel, InvalidJwtToken, err))
	}

	// Check Sign
	jwtToken, err = token.GenerateToken(jwtData.Payload.Purpose, jwtData.Payload.PublicSessionID, jwtData.Payload.IssuedAt)
	if err != nil { // Обработка ошибки
		errLabel = "q2LFd94k"
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", errLabel, InvalidJwtToken, err))
	}

	if strings.Split(jwtToken, ".")[2] != jwtTokenArr[2] {
		errLabel = "xf6qONVc"
		return jwtData, errors.New(fmt.Sprintf("%s:%s", errLabel, InvalidJwtToken))
	}

	// Check Lifetime
	switch jwtData.Payload.Purpose {
	case PropsPurposeAccess:
		lifeTime = token.AccessTokenLifetime
	case PropsPurposeRefresh:
		lifeTime = token.RefreshTokenLifetime
	default:
		errLabel = "n3LDfSbA"
		return jwtData, errors.New(fmt.Sprintf("%s:%s", errLabel, InvalidJwtToken))
	}
	if time.Now().UTC().Unix() >
		time.Unix(jwtData.Payload.IssuedAt, 0).Add(time.Minute*time.Duration(lifeTime)).UTC().Unix() {
		errLabel = "WoZD4wBI"
		return jwtData, errors.New(fmt.Sprintf("%s:%s", errLabel, InvalidJwtToken))
	}

	return jwtData, nil
}
