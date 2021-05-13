/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package jwt

import (
	"github.com/dmalix/financelime-authorization/packages/cryptographer"
	"testing"
	"time"
)

func TestToken_success_secretBoxPlainText(t *testing.T) {

	var err error
	var token string

	jwtManager := NewToken(
		"secretKey",
		ParamSigningAlgorithmHS256,
		"issuer",
		"subject",
		100,
		0)

	token, err = jwtManager.GenerateToken("sessionID", []byte("userData"), ParamPurposeAccess)

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

func TestToken_success_secretBoxEncrypted(t *testing.T) {

	var err error
	var token string
	const dataSource = "data"
	var secretBoxEncrypted []byte

	jwtManager := NewToken(
		"secretKey",
		ParamSigningAlgorithmHS256,
		"issuer",
		"subject",
		100,
		0)

	cryptoManager := cryptographer.NewCryptographer("secretKey")

	secretBoxEncrypted, err = cryptoManager.Encrypt([]byte(dataSource))

	token, err = jwtManager.GenerateToken("sessionID", secretBoxEncrypted, ParamPurposeAccess)

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
		ParamSigningAlgorithmHS256,
		"issuer",
		"subject",
		0,
		0)

	token, err = jwtManager.GenerateToken("sessionID", []byte("data"), ParamPurposeAccess)

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
