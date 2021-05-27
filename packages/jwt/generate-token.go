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
	"fmt"
	"hash"
	"html"
	"time"
)

func (token *Token) GenerateToken(sessionID string, data []byte, tokenPurpose string, issuedAt ...int64) (string, error) {

	var (
		headersBase64 string
		payloadBase64 string
		valueByte     []byte
		err           error
		jwt           string
		signature     string
		unsignedToken string
		mac           hash.Hash
		jsonWebToken  JsonWebToken
	)
	const NoPadding rune = -1

	//    Headers

	jsonWebToken.Headers.Type = ParamTypeJWT
	jsonWebToken.Headers.SigningAlgorithm = token.SigningAlgorithm
	valueByte, err = json.Marshal(jsonWebToken.Headers)
	if err != nil {
		return "", fmt.Errorf("failed convert the jsonWebToken.headers to JSON-format: %s", err)
	}
	headersBase64 = base64.URLEncoding.WithPadding(NoPadding).EncodeToString(valueByte)

	//   Payload

	jsonWebToken.Payload.Issuer = token.Issuer
	jsonWebToken.Payload.Subject = token.Subject

	if tokenPurpose != ParamPurposeAccess && tokenPurpose != ParamPurposeRefresh {
		return "", fmt.Errorf("invalid the tokenPurpose param: %s", err)
	}
	jsonWebToken.Payload.Purpose = tokenPurpose

	jsonWebToken.Payload.PublicSessionID = sessionID

	jsonWebToken.Payload.Data = data
	if len(issuedAt) == 0 {
		jsonWebToken.Payload.IssuedAt = time.Now().UTC().Unix()
	} else {
		jsonWebToken.Payload.IssuedAt = issuedAt[0]
	}

	valueByte, err = json.Marshal(jsonWebToken.Payload)
	if err != nil {
		return "", fmt.Errorf("failed convert the jsonWebToken.Payload to JSON-format: %s", err)
	}
	payloadBase64 = base64.URLEncoding.WithPadding(NoPadding).EncodeToString(valueByte)

	//    Sign

	unsignedToken = headersBase64 + "." + payloadBase64
	switch jsonWebToken.Headers.SigningAlgorithm {
	case ParamSigningAlgorithmHS256:
		mac = hmac.New(sha256.New, []byte(token.SecretKey))
	case ParamSigningAlgorithmHS512:
		mac = hmac.New(sha512.New, []byte(token.SecretKey))
	default:
		return "", fmt.Errorf("invalid algorithm: %s", err)
	}

	mac.Write([]byte(unsignedToken))
	signature = hex.EncodeToString(mac.Sum(nil))

	//   Make JWT

	jwt =
		html.UnescapeString(headersBase64) +
			"." + html.UnescapeString(payloadBase64) +
			"." + html.UnescapeString(signature)

	return jwt, err
}
