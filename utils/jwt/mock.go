package jwt

import "github.com/dmalix/financelime-authorization/models"

type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (s *MockDescription) GenerateToken(publicSessionID string, userData []byte, tokenPurpose string, issuedAt ...int64) (string, error) {
	return "", MockData.Expected.Error
}

func (s *MockDescription) VerifyToken(jwt string) (models.JwtData, error) {
	var data models.JwtData
	return data, MockData.Expected.Error
}
