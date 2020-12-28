package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"io"
)

func (c Cipher) Encrypt(data []byte) ([]byte, error) {

	var (
		err    error
		result []byte

		cipherBlock cipher.Block
		cipherGCM   cipher.AEAD
		nonceUse    []byte
	)

	cipherBlock, err = aes.NewCipher([]byte(createHash(c.SecretKey)))
	if err != nil {
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to create the new AES cipherBlock",
			err))
	}

	cipherGCM, err = cipher.NewGCM(cipherBlock)
	if err != nil {
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to create the new cipherGCM",
			err))
	}

	nonceUse = make([]byte, cipherGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonceUse)
	if err != nil {
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to test for reading full of the nonce used",
			err))
	}

	result = cipherGCM.Seal(nonceUse, nonceUse, data, nil)

	return result, nil
}
