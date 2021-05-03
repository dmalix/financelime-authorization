/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"hash"
	"html"
	"time"
)

func (token *Token) GenerateToken(sessionID string, userData []byte, tokenPurpose string, issuedAt ...int64) (string, error) {

	var (
		headersBase64 string
		payloadBase64 string
		valueByte     []byte
		err           error
		jwt           string
		signature     string
		unsignedToken string
		mac           hash.Hash
		jwtData       JwtData
	)
	const NoPadding rune = -1

	//    Headers
	// ------------

	jwtData.Headers.Type = PropsTypeJWT
	jwtData.Headers.SigningAlgorithm = token.SigningAlgorithm
	valueByte, err = json.Marshal(jwtData.Headers)
	if err != nil {
		return jwt,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed convert the jwtData.headers to JSON-format",
				err))
	}
	headersBase64 = base64.URLEncoding.WithPadding(NoPadding).EncodeToString(valueByte)

	//   Payload
	// ------------

	jwtData.Payload.Issuer = token.Issuer
	jwtData.Payload.Subject = token.Subject

	if tokenPurpose != PropsPurposeAccess && tokenPurpose != PropsPurposeRefresh {
		return jwt,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Invalid the tokenPurpose param",
				err))
	}
	jwtData.Payload.Purpose = tokenPurpose

	jwtData.Payload.PublicSessionID = sessionID

	jwtData.Payload.UserData = userData
	if len(issuedAt) == 0 {
		jwtData.Payload.IssuedAt = time.Now().UTC().Unix()
	} else {
		jwtData.Payload.IssuedAt = issuedAt[0]
	}

	valueByte, err = json.Marshal(jwtData.Payload)
	if err != nil {
		return jwt,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed convert the jwtData.Payload to JSON-format",
				err))
	}
	payloadBase64 = base64.URLEncoding.WithPadding(NoPadding).EncodeToString(valueByte)

	//    Sign
	// -----------

	unsignedToken = headersBase64 + "." + payloadBase64
	switch jwtData.Headers.SigningAlgorithm {
	case PropsSigningAlgorithmHS256:
		mac = hmac.New(sha256.New, []byte(token.SecretKey))
	case PropsSigningAlgorithmHS512:
		mac = hmac.New(sha512.New, []byte(token.SecretKey))
	default:
		return jwt,
			errors.New(fmt.Sprintf("%s:%s",
				trace.GetCurrentPoint(),
				"invalid algorithm"))
	}

	mac.Write([]byte(unsignedToken))
	signature = hex.EncodeToString(mac.Sum(nil))

	//   Make JWT
	// ------------

	jwt =
		html.UnescapeString(headersBase64) +
			"." + html.UnescapeString(payloadBase64) +
			"." + html.UnescapeString(signature)

	return jwt, err
}
