/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"github.com/dmalix/financelime-authorization/models"
	"net/mail"
)

//noinspection GoSnakeCaseUsage
type AddEmailMessageToQueue_MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var AddEmailMessageToQueue_MockData AddEmailMessageToQueue_MockDescription

func (manager AddEmailMessageToQueue_MockDescription) AddEmailMessageToQueue(_ chan models.EmailMessage,
	to mail.Address, _, _ string, _ ...string) error {
	return AddEmailMessageToQueue_MockData.Expected.Error
}
