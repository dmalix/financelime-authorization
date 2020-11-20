package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func (c Cipher) Encrypt(data []byte) ([]byte, error) {

	var (
		err      error
		errLabel string
		result   []byte

		cipherBlock cipher.Block
		cipherGCM   cipher.AEAD
		nonceUse    []byte
	)

	cipherBlock, err = aes.NewCipher([]byte(createHash(c.SecretKey)))
	if err != nil {
		errLabel = "a49WZPGR"
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to create the new AES cipherBlock",
			err))
	}

	cipherGCM, err = cipher.NewGCM(cipherBlock)
	if err != nil {
		errLabel = "HVQKTs7X"
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to create the new cipherGCM",
			err))
	}

	nonceUse = make([]byte, cipherGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonceUse)
	if err != nil {
		errLabel = "QdBJojn3"
		return result, errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to test for reading full of the nonce used",
			err))
	}

	result = cipherGCM.Seal(nonceUse, nonceUse, data, nil)

	return result, nil
}
