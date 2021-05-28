/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import "net/mail"

func (manager Manager) AddEmailMessageToQueue(messageQueue chan MessageBox, request Request, email Email) error {

	var messageBox MessageBox
	var from = mail.Address{}

	messageBox.Email.To = email.To
	if email.From == from {
		messageBox.Email.From = manager.from
	}
	messageBox.Email.Subject = email.Subject
	messageBox.Email.Body = email.Body
	messageBox.Email.MessageID = email.MessageID

	messageBox.Request.RemoteAddr = request.RemoteAddr
	messageBox.Request.RemoteAddrKey = request.RemoteAddrKey
	messageBox.Request.RequestID = request.RequestID
	messageBox.Request.RequestIDKey = request.RequestIDKey

	messageQueue <- messageBox

	return nil
}
