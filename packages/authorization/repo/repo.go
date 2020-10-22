/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repo

import (
	"database/sql"
)

type Repo struct {
	dbAuthMain *sql.DB
	dbAuthRead *sql.DB
	dbBlade    *sql.DB
}

func NewAuthorizationRepo(dbAuthMain, dbAuthRead, dbBlade *sql.DB) *Repo {
	return &Repo{
		dbAuthMain: dbAuthMain,
		dbAuthRead: dbAuthRead,
		dbBlade:    dbBlade,
	}
}
