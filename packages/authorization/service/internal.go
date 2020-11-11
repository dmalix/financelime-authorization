package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/utils/random"
	"hash"
	"strconv"
	"time"
)

func (s *Service) generatePublicID(privateID int64) (string, error) {

	var (
		err             error
		errLabel        string
		hs              hash.Hash
		publicSessionID string
	)

	hs = sha256.New()
	_, err = hs.Write([]byte(
		strconv.FormatInt(privateID, 10) +
			random.StringRand(16, 16, false) +
			time.Now().String() +
			s.config.CryptoSalt))
	if err != nil {
		errLabel = "7ZHDXTE3"
		return publicSessionID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	publicSessionID = hex.EncodeToString(hs.Sum(nil))

	return publicSessionID, nil
}
