/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
)

func (daemon SenderDeamon) smtpSender(to, from mail.Address, subject, body string, messageID string) error {

	var (
		smtpAuth       smtp.Auth
		smtpTlsConfig  *tls.Config
		smtpTlsConnect *tls.Conn
		smtpClient     *smtp.Client
		smtpMessage    io.WriteCloser

		messageHeaders map[string]string
		messageResult  string

		err error
	)

	// 1. Preparing the messageHeaders and body of the message

	messageHeaders = make(map[string]string)

	if len(messageID) > 0 {
		messageHeaders["Message-Id"] = messageID
	}

	messageHeaders["From"] = from.String()
	messageHeaders["To"] = to.String()
	messageHeaders["Subject"] = subject
	messageHeaders["MIME-Version"] = "1.0"
	messageHeaders["Content-Type"] = "text/plain; charset=\"utf-8\""
	messageHeaders["Content-Transfer-Encoding"] = "base64"

	messageResult = ""

	for k, v := range messageHeaders {
		messageResult += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	messageResult += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// 2. Prepare authorization data

	smtpAuth = smtp.PlainAuth("", daemon.AuthSMTPUser, daemon.AuthSMTPPassword, daemon.AuthSMTPHost)

	// 3. Performing a TLS Connection

	smtpTlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         daemon.AuthSMTPHost,
	}

	smtpTlsConnect, err =
		tls.Dial("tcp", fmt.Sprintf("%s:%d", daemon.AuthSMTPHost, daemon.AuthSMTPPort), smtpTlsConfig)
	if err != nil {
		return fmt.Errorf("failed to performing a TLS connection: %s", err)
	}

	// 4. Connect client

	smtpClient, err =
		smtp.NewClient(smtpTlsConnect, daemon.AuthSMTPHost)
	if err != nil {
		return fmt.Errorf("failed to Perform client authentication: %s", err)
	}

	// 5. Perform client authentication

	err = smtpClient.Auth(smtpAuth)
	if err != nil {
		return fmt.Errorf("failed to Perform client authentication: %s", err)
	}

	// 6. Send the data of the sender and recipient

	err = smtpClient.Mail(from.Address)
	if err != nil {
		return fmt.Errorf("failed to Send the data of the sender: %s", err)
	}

	err = smtpClient.Rcpt(to.Address)
	if err != nil {
		return fmt.Errorf("failed to Send the data of the recipient: %s", err)
	}

	// 7. Send message

	smtpMessage, err = smtpClient.Data()
	if err != nil {
		return fmt.Errorf("failed to Send message (DATA): %s", err)
	}

	_, err = smtpMessage.Write([]byte(messageResult))
	if err != nil {
		return fmt.Errorf("failed to Send message (WRITE): %s", err)
	}

	err = smtpMessage.Close()
	if err != nil {
		return fmt.Errorf("failed to Send message (CLOSE): %s", err)
	}

	// 8. Quit

	err = smtpClient.Quit()
	if err != nil {
		return fmt.Errorf("failed to Quit: %s", err)
	}

	return nil
}
