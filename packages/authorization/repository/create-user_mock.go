/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"github.com/dmalix/financelime-rest-api/models"
)

//noinspection GoSnakeCaseUsage
type CreateUser_MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var CreateUser_MockData CreateUser_MockDescription

func (repo *CreateUser_MockDescription) CreateUser(_ *models.User, _, _ string, _ bool) error {

	if CreateUser_MockData.Expected.Error != nil {
		return CreateUser_MockData.Expected.Error
	}

	return CreateUser_MockData.Expected.Error
}
