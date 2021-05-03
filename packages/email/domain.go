package email

import (
	"net/mail"
)

type Message interface {
	AddEmailMessageToQueue(messageQueue chan EmailMessage, to mail.Address, subject, body string, messageID ...string) error
}
