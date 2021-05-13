/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

type JsonWebToken struct {
	Headers struct {
		SigningAlgorithm string `json:"alg"`
		Type             string `json:"typ"`
	}
	Payload struct {
		Issuer          string `json:"iss"`
		Subject         string `json:"sub"`
		Purpose         string `json:"purpose"`
		PublicSessionID string `json:"sessionID"`
		Data            []byte `json:"data"`
		IssuedAt        int64  `json:"iat"`
	}
}
