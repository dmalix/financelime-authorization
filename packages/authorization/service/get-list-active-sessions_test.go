/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"testing"
)

func TestGetListActiveSessions(t *testing.T) {

	var (
		encryptedUserData []byte
		service           = new(MockDescription)
	)

	_, err := service.GetListActiveSessions(encryptedUserData)

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}
