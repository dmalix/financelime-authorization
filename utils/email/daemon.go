/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"github.com/dmalix/financelime-rest-api/models"
	"log"
	"time"
)

func (daemon Daemon) Run() {

	var message models.EmailMessage
	var err error

	for {

		select {

		case message = <-daemon.MessageQueue:

			for {
				err = daemon.smtpSender(message.To, message.From, message.Subject, message.Body, message.MessageID)
				if err == nil {
					break
				}

				log.Printf("FATAL [%s: Failed to send a email message [%s]]", "a3GXhiMR", err)

				time.Sleep(time.Second * time.Duration(5))
			}

			time.Sleep(time.Second * time.Duration(3))

		default:

		}
	}
}
