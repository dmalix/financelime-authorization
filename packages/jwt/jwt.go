/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

type Token struct {
	SecretKey               string
	SigningAlgorithm        string
	Issuer                  string
	Subject                 string
	AccessTokenLifetimeSec  int
	RefreshTokenLifetimeSec int
}

func NewToken(
	SecretKey string,
	SigningAlgorithm string,
	Issuer string,
	Subject string,
	AccessTokenLifetimeSec int,
	RefreshTokenLifetimeSec int) *Token {
	return &Token{
		SecretKey:               SecretKey,
		SigningAlgorithm:        SigningAlgorithm,
		Issuer:                  Issuer,
		Subject:                 Subject,
		AccessTokenLifetimeSec:  AccessTokenLifetimeSec,
		RefreshTokenLifetimeSec: RefreshTokenLifetimeSec,
	}
}

const (
	ParamTypeJWT               = "JWT"
	ParamPurposeAccess         = "access"
	ParamPurposeRefresh        = "refresh"
	ParamSigningAlgorithmHS256 = "HS256"
	ParamSigningAlgorithmHS512 = "HS512"
)