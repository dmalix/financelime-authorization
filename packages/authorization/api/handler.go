/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"github.com/dmalix/financelime-authorization/packages/authorization"
)

type Handler struct {
	service authorization.Service
}

func NewHandler(service authorization.Service) *Handler {
	return &Handler{
		service: service,
	}
}

const (
	ContextPublicSessionID     = "publicSessionID"
	ContextEncryptedUserData   = "encryptedUserData"
	ContentTypeApplicationJson = "application/json;charset=utf-8"
	ContentTypeTextPlain       = "text/plain;charset=utf-8"
)
