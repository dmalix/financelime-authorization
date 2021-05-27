package middleware

import (
	"context"
	"fmt"
)

type contextGetter struct {
}

func NewContextGetter() *contextGetter {
	return &contextGetter{}
}

func (cg *contextGetter) GetRemoteAddr(ctx context.Context) (string, error) {
	return getRemoteAddr(ctx)
}

func getRemoteAddr(ctx context.Context) (string, error) {
	ctxValue := ctx.Value(ContextKeyRemoteAddr)
	if ctxValue == nil {
		return "", fmt.Errorf("failed to get context %s", ContextKeyRemoteAddr)
	}
	return ctxValue.(string), nil
}

func (cg *contextGetter) GetRequestID(ctx context.Context) (string, string, error) {
	return getRequestID(ctx)
}

func getRequestID(ctx context.Context) (string, string, error) {
	ctxValue := ctx.Value(ContextKeyRequestID)
	if ctxValue == nil {
		return "", "", fmt.Errorf("failed to get context %s", ContextKeyRequestID)
	}
	return ctxValue.(string), ContextKeyRequestID, nil
}

func (cg *contextGetter) GetPublicSessionID(ctx context.Context) (string, error) {
	ctxValue := ctx.Value(ContextKeyPublicSessionID)
	if ctxValue == nil {
		return "", fmt.Errorf("failed to get context %s", ContextKeyPublicSessionID)
	}
	return ctxValue.(string), nil
}

func (cg *contextGetter) GetEncryptedUserData(ctx context.Context) ([]byte, error) {
	ctxValue := ctx.Value(ContextKeyEncryptedJWTData)
	if ctxValue == nil {
		return nil, fmt.Errorf("failed to get context %s", ContextKeyEncryptedJWTData)
	}
	return ctxValue.([]byte), nil
}
