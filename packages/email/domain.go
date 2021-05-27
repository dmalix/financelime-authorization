package email

import (
	"context"
	"go.uber.org/zap"
)

type Message interface {
	AddEmailMessageToQueue(messageQueue chan MessageBox, request Request, email Email) error
}

type Daemon interface {
	Run(ctx context.Context, logger *zap.Logger)
}
