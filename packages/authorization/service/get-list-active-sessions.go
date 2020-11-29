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
	   	Get a list of active sessions
	   		----------------
	   		Return:
				sessions []models.Session
	   			err error
*/
// Related interfaces:
//	packages/authorization/domain.go
func (s *Service) GetListActiveSessions(encryptedUserData []byte) ([]models.Session, error) {

	var (
		err               error
		errLabel          string
		sessions          []models.Session
		decryptedUserData []byte
		user              models.User
	)

	decryptedUserData, err = s.cryptographer.Decrypt([]byte(encryptedUserData))
	if err != nil {
		return sessions,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to decrypt the user data",
				err))
	}

	err = json.Unmarshal(decryptedUserData, &user)
	if err != nil {
		errLabel = "9AZS3RH1"
		return sessions,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"Failed to unmarshal the decryptedUserData value to struct [%s]",
				err))
	}

	sessions, err = s.repository.GetListActiveSessions(user.ID)
	if err != nil {
		errLabel = "3RAZSH91"
		return sessions,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				errLabel,
				"a system error was returned",
				err))
	}
	return sessions, nil
}
