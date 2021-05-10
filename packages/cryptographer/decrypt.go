package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
)

func (c Cipher) Decrypt(data []byte) ([]byte, error) {

	var (
		err    error
		result []byte

		hashedKey   []byte
		cipherBlock cipher.Block
		cipherGCM   cipher.AEAD
		nonceUse    []byte
	)

	hashedKey = []byte(createHash(c.SecretKey))
	cipherBlock, err = aes.NewCipher(hashedKey)
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

	nonceSize := cipherGCM.NonceSize()
	nonceUse, ciphertext := data[:nonceSize], data[nonceSize:]
	result, err = cipherGCM.Open(nil, nonceUse, ciphertext, nil)
	if err != nil {
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to open decrypts",
			err))
	}
	return result, nil
}
