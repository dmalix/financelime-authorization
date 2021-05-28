/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"net/mail"
)

type SenderDeamon struct {
	AuthSMTPUser     string
	AuthSMTPPassword string
	AuthSMTPHost     string
	AuthSMTPPort     int
	MessageQueue     chan MessageBox
}

func NewSenderDaemon(
	authSMTPUser,
	authSMTPPassword,
	authSMTPHost string,
	authSMTPPort int,
	messageQueue chan MessageBox) *SenderDeamon {
	return &SenderDeamon{
		AuthSMTPUser:     authSMTPUser,
		AuthSMTPPassword: authSMTPPassword,
		AuthSMTPHost:     authSMTPHost,
		AuthSMTPPort:     authSMTPPort,
		MessageQueue:     messageQueue,
	}
}

type Manager struct {
	from mail.Address
}

func NewManager(
	from mail.Address) *Manager {
	return &Manager{
		from: from,
	}
}
