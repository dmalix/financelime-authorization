/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package middleware

import (
	"github.com/dmalix/financelime-authorization/packages/jwt"
)

type ConfigMiddleware struct {
	RequestIDRequired bool
	RequestIDCheck    bool
}

type Middleware struct {
	config ConfigMiddleware
	jwt    jwt.Jwt
}

func NewMiddleware(
	config ConfigMiddleware,
	jwt jwt.Jwt) *Middleware {
	return &Middleware{
		config: config,
		jwt:    jwt,
	}
}
