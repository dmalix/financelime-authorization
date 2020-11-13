/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package api

import (
	"github.com/dmalix/financelime-authorization/packages/system"
)

type Handler struct {
	service system.Service
}

func NewHandler(service system.Service) *Handler {
	return &Handler{
		service: service,
	}
}
