package information

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type Rest interface {
	Version(logger *zap.Logger) http.Handler
}

type Service interface {
	Version(ctx context.Context, logger *zap.Logger) (string, string, error)
}
