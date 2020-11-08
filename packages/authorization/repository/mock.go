/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"github.com/dmalix/financelime-rest-api/models"
)

//noinspection GoSnakeCaseUsage
type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (repo *MockDescription) CreateUser(_ *models.User, _, _, _ string, _ bool) error {
	return MockData.Expected.Error
}

func (repo *MockDescription) ConfirmUserEmail(_ string) (models.User, error) {
	var user models.User
	return user, MockData.Expected.Error
}
