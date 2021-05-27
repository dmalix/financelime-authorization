/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func (token *Token) VerifyToken(jwt string) (JsonWebToken, error) {

	var (
		valueByte    []byte
		err          error
		jwtTokenArr  []string
		jwtToken     string
		lifeTime     int
		jsonWebToken JsonWebToken
	)

	const messageInvalidJwtToken = "invalid JWT-token: %s"
	const NoPadding rune = -1

	jwtTokenArr = strings.Split(jwt, ".")

	// Check Headers
	valueByte, err = base64.URLEncoding.WithPadding(NoPadding).DecodeString(jwtTokenArr[0])
	if err != nil {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}

	err = json.Unmarshal(valueByte, &jsonWebToken.Headers)
	if err != nil {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}

	// Check Payload
	valueByte, err = base64.URLEncoding.WithPadding(NoPadding).DecodeString(jwtTokenArr[1])
	if err != nil {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}
	err = json.Unmarshal(valueByte, &jsonWebToken.Payload)
	if err != nil {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}

	// Check Sign
	jwtToken, err = token.GenerateToken(jsonWebToken.Payload.PublicSessionID, jsonWebToken.Payload.Data,
		jsonWebToken.Payload.Purpose, jsonWebToken.Payload.IssuedAt)
	if err != nil {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}

	if strings.Split(jwtToken, ".")[2] != jwtTokenArr[2] {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}

	// Check Lifetime
	switch jsonWebToken.Payload.Purpose {
	case ParamPurposeAccess:
		lifeTime = token.AccessTokenLifetimeSec
	case ParamPurposeRefresh:
		lifeTime = token.RefreshTokenLifetimeSec
	default:
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}
	if time.Now().UTC().Unix() >
		time.Unix(jsonWebToken.Payload.IssuedAt, 0).Add(time.Second*time.Duration(lifeTime)).UTC().Unix() {
		return JsonWebToken{}, fmt.Errorf(messageInvalidJwtToken, err)
	}

	return jsonWebToken, nil
}
