package jwt

type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

var MockData MockDescription

func (s *MockDescription) GenerateToken(publicSessionID string, data []byte, tokenPurpose string, issuedAt ...int64) (string, error) {
	return "", MockData.Expected.Error
}

func (s *MockDescription) VerifyToken(jwt string) (JsonWebToken, error) {
	var data JsonWebToken
	return data, MockData.Expected.Error
}
