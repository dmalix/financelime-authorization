package cryptographer

import (
	"encoding/json"
)

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
	result := map[string]interface{}{
		"test_data": "test_data",
	}
	bytesRepresentation, _ := json.Marshal(result)
	return bytesRepresentation, MockData.Expected.Error
}
