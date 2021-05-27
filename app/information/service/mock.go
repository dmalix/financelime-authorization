package service

import (
	"context"
	"go.uber.org/zap"
)

type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

var MockData MockDescription

func (a *MockDescription) Version(_ context.Context, _ *zap.Logger) (string, string, error) {
	return "version", "build", MockData.Expected.Error
}
