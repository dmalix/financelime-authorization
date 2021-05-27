package middleware

import "context"

type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

func (s *MockDescription) GetRemoteAddr(ctx context.Context) (string, error) {
	return "", nil
}
func (s *MockDescription) GetRequestID(ctx context.Context) (string, string, error) {
	return "", "", nil
}
func (s *MockDescription) GetPublicSessionID(ctx context.Context) (string, error) {
	return "", nil
}
func (s *MockDescription) GetEncryptedUserData(ctx context.Context) ([]byte, error) {
	return nil, nil
}
