/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"context"
	"go.uber.org/zap"
	"time"
)

func (daemon SenderDeamon) Run(ctx context.Context, logger *zap.Logger) {

	var message EMessage
	var err error

	for {

		select {

		case <-ctx.Done():
			return

		case message = <-daemon.MessageQueue:
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						err = daemon.smtpSender(message.To, message.From, message.Subject, message.Body, message.MessageID)
						if err == nil {
							return
						}
						logger.DPanic("failed to send a email message", zap.Error(err))
						time.Sleep(time.Second * time.Duration(3))
					}
				}
			}()
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}
