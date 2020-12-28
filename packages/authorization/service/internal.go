package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/random"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"hash"
	"strconv"
	"time"
)

func (s *Service) generatePublicID(privateID int64) (string, error) {

	var (
		err             error
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
		return publicSessionID, errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	publicSessionID = hex.EncodeToString(hs.Sum(nil))

	return publicSessionID, nil
}
