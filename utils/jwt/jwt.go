/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

type Token struct {
	SecretKey            string
	SigningAlgorithm     string
	Issuer               string
	Subject              string
	AccessTokenLifetime  int
	RefreshTokenLifetime int
}

func NewToken(
	SecretKey string,
	SigningAlgorithm string,
	Issuer string,
	Subject string,
	AccessTokenLifetime int,
	RefreshTokenLifetime int) *Token {
	return &Token{
		SecretKey,
		SigningAlgorithm,
		Issuer,
		Subject,
		AccessTokenLifetime,
		RefreshTokenLifetime,
	}
}

const (
	PropsTypeJWT               = "JWT"
	PropsPurposeAccess         = "access"
	PropsPurposeRefresh        = "refresh"
	PropsSigningAlgorithmHS256 = "HS256"
	PropsSigningAlgorithmHS512 = "HS512"
)
