/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"github.com/dmalix/financelime-rest-api/models"
)

type Daemon struct {
	AuthSMTPUser     string
	AuthSMTPPassword string
	AuthSMTPHost     string
	AuthSMTPPort     string
	MessageQueue     chan models.EmailMessage
}

type Manager struct {
}

func NewSenderDaemon(
	authSMTPUser,
	authSMTPPassword,
	authSMTPHost,
	authSMTPPort string,
	messageQueue chan models.EmailMessage) *Daemon {
	return &Daemon{
		AuthSMTPUser:     authSMTPUser,
		AuthSMTPPassword: authSMTPPassword,
		AuthSMTPHost:     authSMTPHost,
		AuthSMTPPort:     authSMTPPort,
		MessageQueue:     messageQueue,
	}
}

func NewManager() *Manager {
	return &Manager{}
}
