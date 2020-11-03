/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"github.com/dmalix/financelime-rest-api/models"
	"net/mail"
)

func (manager Manager) AddEmailMessageToQueue(messageQueue chan models.EmailMessage,
	to mail.Address, subject, body string, messageID ...string) error {

	var message models.EmailMessage

	message.To = to
	message.From = manager.from
	message.Subject = subject
	message.Body = body
	if len(message.MessageID) > 0 {
		message.MessageID = messageID[0]
	}

	messageQueue <- message

	return nil
}
