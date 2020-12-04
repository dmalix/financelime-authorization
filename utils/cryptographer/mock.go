package cryptographer

type MockDescription struct {
	Props    struct{}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (s *MockDescription) Encrypt(data []byte) ([]byte, error) {
	return nil, MockData.Expected.Error
}

func (s *MockDescription) Decrypt(data []byte) ([]byte, error) {
	return nil, MockData.Expected.Error
}
