/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

import (
	"github.com/dmalix/financelime-authorization/utils/cryptographer"
	"testing"
	"time"
)

func TestToken_success_userDataPlainText(t *testing.T) {

	var err error
	var token string

	jwtManager := NewToken(
		"secretKey",
		PropsSigningAlgorithmHS256,
		"issuer",
		"subject",
		100,
		0)

	token, err = jwtManager.GenerateToken("sessionID", []byte("userData"), PropsPurposeAccess)

	if err != nil {
		t.Errorf("function returned wrong error value: got %v want %v",
			err, nil)
	}

	_, err = jwtManager.VerifyToken(token)

	if err != nil {
		t.Errorf("function returned wrong error value: got %v want %v",
			err, nil)
	}
}

func TestToken_success_userDataEncrypted(t *testing.T) {

	var err error
	var token string
	const userDataSource = "userData"
	var userDataEncrypted []byte

	jwtManager := NewToken(
		"secretKey",
		PropsSigningAlgorithmHS256,
		"issuer",
		"subject",
		100,
		0)

	cryptoManager := cryptographer.NewCryptographer("secretKey")

	userDataEncrypted, err = cryptoManager.Encrypt([]byte(userDataSource))

	token, err = jwtManager.GenerateToken("sessionID", userDataEncrypted, PropsPurposeAccess)

	if err != nil {
		t.Errorf("function returned wrong error value: got %v want %v",
			err, nil)
	}

	_, err = jwtManager.VerifyToken(token)

	if err != nil {
		t.Errorf("function returned wrong error value: got %v want %v",
			err, nil)
	}
}

func TestToken_invalid(t *testing.T) {

	var err error
	var token string

	jwtManager := NewToken(
		"secretKey",
		PropsSigningAlgorithmHS256,
		"issuer",
		"subject",
		0,
		0)

	token, err = jwtManager.GenerateToken("sessionID", []byte("userData"), PropsPurposeAccess)

	if err != nil {
		t.Errorf("function returned wrong error value: got %v want %v",
			err, nil)
	}

	time.Sleep(1 * time.Second)

	_, err = jwtManager.VerifyToken(token)

	if err == nil {
		t.Errorf("function returned wrong error value: got %v want %v",
			err, "!nil")
	}
}
