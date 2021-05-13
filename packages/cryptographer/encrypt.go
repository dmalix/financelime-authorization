package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
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
		return result, fmt.Errorf("failed to create the new AES cipherBlock: %s", err)
	}

	cipherGCM, err = cipher.NewGCM(cipherBlock)
	if err != nil {
		return result, fmt.Errorf("failed to create the new cipherGCM: %s", err)
	}

	nonceUse = make([]byte, cipherGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonceUse)
	if err != nil {
		return result, fmt.Errorf("failed to test for reading full of the nonce used: %s", err)
	}

	result = cipherGCM.Seal(nonceUse, nonceUse, data, nil)

	return result, nil
}
