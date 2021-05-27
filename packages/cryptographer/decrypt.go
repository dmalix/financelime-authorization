package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
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
		return result, fmt.Errorf("failed to create the new AES cipherBlock: %s", err)
	}

	cipherGCM, err = cipher.NewGCM(cipherBlock)
	if err != nil {
		return result, fmt.Errorf("failed to create the new cipherGCM: %s", err)
	}

	nonceSize := cipherGCM.NonceSize()
	nonceUse, ciphertext := data[:nonceSize], data[nonceSize:]
	result, err = cipherGCM.Open(nil, nonceUse, ciphertext, nil)
	if err != nil {
		return result, fmt.Errorf("failed to open decrypts: %s", err)
	}
	return result, nil
}
