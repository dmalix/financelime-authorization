/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"github.com/dmalix/financelime-rest-api/models"
	"net/mail"
)

type UserMessage interface {
	AddEmailMessageToQueue(messageQueue chan models.EmailMessage, to mail.Address, subject, body string, messageID ...string) error
}
