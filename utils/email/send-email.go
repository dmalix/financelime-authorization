/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
)

/*
	Send Email
		------------------
			Return:
				error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details])::
				        ------------------------------------------------
				        TLS_CONNECTION:  failed to Performing a TLS Connection
				        CONNECT_CLIENT:  failed to Connect New Client
				        CLIENT_AUTH:     failed to Perform client authentication
				        DATA_SENDER:     failed to Send the data of the sender
				        DATA_RECIPIENT:  failed to Send the data of the recipient
				        MESSAGE_DATA:    failed to Send message (DATA)
				        MESSAGE_WRITE:   failed to Send message (WRITE)
				        MESSAGE_CLOSE:   failed to Send message (CLOSE)
				        QUIT:            failed to Quit
*/
// Related interfaces:
//	packages/authorization/domain/user.go
func (authSMTP *AuthSMTP) SendEmail(to mail.Address, subject, body string, messageID ...string) error {

	var from = mail.Address{Name: "SUPPORT_NAME", Address: "support@financelime.com"}

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
	// -------------------------------------------------------

	messageHeaders = make(map[string]string)

	if len(messageID) == 1 {
		messageHeaders["Message-Id"] = messageID[0]
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
	// -----------------------------

	smtpAuth = smtp.PlainAuth("", authSMTP.User, authSMTP.Password, authSMTP.Host)

	// 3. Performing a TLS Connection
	// ------------------------------

	// TLS configuration

	smtpTlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         authSMTP.Host,
	}

	smtpTlsConnect, err =
		tls.Dial("tcp", fmt.Sprintf("%s:%s", authSMTP.Host, authSMTP.Port), smtpTlsConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"TLS_CONNECTION", "Failed to Performing a TLS Connection", err))
	}

	// 4. Connect client
	// -----------------

	smtpClient, err =
		smtp.NewClient(smtpTlsConnect, authSMTP.Host)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"CONNECT_CLIENT", "failed to Perform client authentication", err))
	}

	// 5. Perform client authentication
	// --------------------------------

	err = smtpClient.Auth(smtpAuth)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"CLIENT_AUTH", "failed to Perform client authentication", err))
	}

	// 6. Send the data of the sender and recipient
	// --------------------------------------------

	err = smtpClient.Mail(from.Address)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"DATA_SENDER", "failed to Send the data of the sender", err))
	}

	err = smtpClient.Rcpt(to.Address)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"DATA_RECIPIENT", "failed to Send the data of the recipient", err))
	}

	// 7. Send message
	// ---------------

	smtpMessage, err = smtpClient.Data()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"MESSAGE_DATA", "failed to Send message (DATA)", err))
	}

	_, err = smtpMessage.Write([]byte(messageResult))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"MESSAGE_WRITE", "failed to Send message (WRITE)", err))
	}

	err = smtpMessage.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"MESSAGE_CLOSE", "failed to Send message (CLOSE)", err))
	}

	// 8. Quit
	// -------

	err = smtpClient.Quit()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"QUIT", "failed to Quit", err))
	}

	return nil

}
