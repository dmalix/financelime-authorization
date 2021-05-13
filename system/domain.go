package system

import (
	"context"
	"net/http"
)

type API interface {
	version(ctx context.Context) http.Handler
}

type Service interface {
	version(ctx context.Context) (string, string, error)
}
