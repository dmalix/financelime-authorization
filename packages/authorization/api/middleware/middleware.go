/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package middleware

import (
	"github.com/dmalix/financelime-authorization/packages/authorization"
)

type Config struct {
	RequestIDRequired bool
	RequestIDCheck    bool
}

type Middleware struct {
	config Config
	jwt    authorization.Jwt
}

func NewMiddleware(
	config Config,
	jwt authorization.Jwt) *Middleware {
	return &Middleware{
		config: config,
		jwt:    jwt,
	}
}
