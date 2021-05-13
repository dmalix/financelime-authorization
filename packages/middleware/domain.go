/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type Middleware interface {
	RequestID(logger *zap.Logger) func(http.Handler) http.Handler
	RemoteAddr(logger *zap.Logger) func(http.Handler) http.Handler
	Logging(logger *zap.Logger) func(http.Handler) http.Handler
	Authorization(logger *zap.Logger) func(http.Handler) http.Handler
}

type ContextGetter interface {
	GetRemoteAddr(ctx context.Context) (string, error)
	GetRequestID(ctx context.Context) (string, string, error)
	GetPublicSessionID(ctx context.Context) (string, error)
	GetEncryptedUserData(ctx context.Context) ([]byte, error)
}
