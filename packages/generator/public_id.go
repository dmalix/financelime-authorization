package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func GeneratePublicID(privateID int64) (string, error) {

	hs := sha256.New()

	_, err := hs.Write([]byte(
		strconv.FormatInt(privateID, 10) + StringRand(16, 16, false) + time.Now().String()))
	if err != nil {
		return "", fmt.Errorf("failed to generate a hash :%s", err)
	}

	publicSessionID := hex.EncodeToString(hs.Sum(nil))

	return publicSessionID, nil
}
