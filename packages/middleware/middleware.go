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

type middleware struct {
	config ConfigMiddleware
	jwt    jwt.JWT
}

func NewMiddleware(
	config ConfigMiddleware,
	jwt jwt.JWT) *middleware {
	return &middleware{
		config: config,
		jwt:    jwt,
	}
}
