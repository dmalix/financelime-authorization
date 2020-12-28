/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"strings"
	"time"
)

func (token *Token) VerifyToken(jwt string) (models.JwtData, error) {

	var (
		valueByte   []byte
		err         error
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
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", trace.GetCurrentPoint(), InvalidJwtToken, err))
	}

	err = json.Unmarshal(valueByte, &jwtData.Headers)
	if err != nil {
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", trace.GetCurrentPoint(), InvalidJwtToken, err))
	}

	// Check Payload
	valueByte, err = base64.URLEncoding.WithPadding(NoPadding).DecodeString(jwtTokenArr[1])
	if err != nil {
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", trace.GetCurrentPoint(), InvalidJwtToken, err))
	}
	err = json.Unmarshal(valueByte, &jwtData.Payload)
	if err != nil {
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", trace.GetCurrentPoint(), InvalidJwtToken, err))
	}

	// Check Sign
	jwtToken, err = token.GenerateToken(jwtData.Payload.PublicSessionID, jwtData.Payload.UserData,
		jwtData.Payload.Purpose, jwtData.Payload.IssuedAt)
	if err != nil { // Обработка ошибки
		return jwtData, errors.New(fmt.Sprintf("%s:%s[%s]", trace.GetCurrentPoint(), InvalidJwtToken, err))
	}

	if strings.Split(jwtToken, ".")[2] != jwtTokenArr[2] {
		return jwtData, errors.New(fmt.Sprintf("%s:%s", trace.GetCurrentPoint(), InvalidJwtToken))
	}

	// Check Lifetime
	switch jwtData.Payload.Purpose {
	case PropsPurposeAccess:
		lifeTime = token.AccessTokenLifetimeSec
	case PropsPurposeRefresh:
		lifeTime = token.RefreshTokenLifetimeSec
	default:
		return jwtData, errors.New(fmt.Sprintf("%s:%s", trace.GetCurrentPoint(), InvalidJwtToken))
	}
	if time.Now().UTC().Unix() >
		time.Unix(jwtData.Payload.IssuedAt, 0).Add(time.Second*time.Duration(lifeTime)).UTC().Unix() {
		return jwtData, errors.New(fmt.Sprintf("%s:%s", trace.GetCurrentPoint(), InvalidJwtToken))
	}

	return jwtData, nil
}
