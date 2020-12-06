/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
)

/*
	Revoke Access Token
		----------------
		Return:
			err             error  - system error
*/
func (s *Service) RevokeRefreshToken(encryptedUserData []byte, publicSessionID string) error {

	var (
		err      error
		errLabel string

		user              models.User
		decryptedUserData []byte
	)

	decryptedUserData, err = s.cryptographer.Decrypt(encryptedUserData)
	if err != nil {
		errLabel = "1S3RH9AZ"
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to decrypt the user data",
			err))
	}

	err = json.Unmarshal(decryptedUserData, &user)
	if err != nil {
		errLabel = "AZRH9S31"
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to unmarshal the decryptedUserData value to struct [%s]",
			err))
	}

	err = s.repository.DeleteSession(user.ID, publicSessionID)
	if err != nil {
		errLabel = "SAZRH931"
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			errLabel,
			"Failed to detele the session [%s]",
			err))
	}

	return nil
}
