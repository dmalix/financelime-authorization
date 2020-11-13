package jwt

import "github.com/dmalix/financelime-authorization/models"

type MockDescription struct {
	Props struct {
		Email      string
		InviteCode string
		Language   string
		RemoteAddr string
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (s *MockDescription) GenerateToken(_, _ string, _ ...int64) (string, error) {
	return "", MockData.Expected.Error
}

func (s *MockDescription) VerifyToken(_ string) (models.JwtData, error) {
	var data models.JwtData
	return data, MockData.Expected.Error
}
