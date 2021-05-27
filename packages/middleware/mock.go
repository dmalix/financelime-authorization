package middleware

import "context"

type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

func (s *MockDescription) GetRemoteAddr(_ context.Context) (string, string, error) {
	return "", "", nil
}
func (s *MockDescription) GetRequestID(_ context.Context) (string, string, error) {
	return "", "", nil
}
func (s *MockDescription) GetPublicSessionID(_ context.Context) (string, error) {
	return "", nil
}
func (s *MockDescription) GetEncryptedUserData(_ context.Context) ([]byte, error) {
	return nil, nil
}
