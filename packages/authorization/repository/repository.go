/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"database/sql"
)

type ConfigRepository struct {
	CryptoSalt              string
	JwtRefreshTokenLifetime int
}

type Repository struct {
	config     ConfigRepository
	dbAuthMain *sql.DB
	dbAuthRead *sql.DB
	dbBlade    *sql.DB
}

func NewRepository(
	config ConfigRepository,
	dbAuthMain,
	dbAuthRead,
	dbBlade *sql.DB) *Repository {
	return &Repository{
		config:     config,
		dbAuthMain: dbAuthMain,
		dbAuthRead: dbAuthRead,
		dbBlade:    dbBlade,
	}
}
