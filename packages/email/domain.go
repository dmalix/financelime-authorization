package email

import (
	"context"
	"go.uber.org/zap"
	"net/mail"
)

type Message interface {
	AddEmailMessageToQueue(messageQueue chan EMessage, to mail.Address, subject, body string, messageID ...string) error
}

type Daemon interface {
	Run(ctx context.Context, logger *zap.Logger)
}
