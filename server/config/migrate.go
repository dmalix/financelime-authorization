/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package config

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
)

func Migrate(db *sql.DB, dropFile, createFile, insertFile string) error {

	var (
		err  error
		body []byte
	)

	if len(dropFile) != 0 {
		body, err = ioutil.ReadFile(dropFile)
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"F1eut9GP",
				"Failed to read the 'drop' file",
				err.Error()))
		}
		_, err = db.Exec(string(body))
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"pPb8wN3J",
				"Failed to complete the request",
				err.Error()))
		}
	}

	if len(createFile) != 0 {
		body, err = ioutil.ReadFile(createFile)
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"VrftP0sD",
				"Failed to read the 'create' file",
				err.Error()))
		}
		_, err = db.Exec(string(body))
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"ktaE1fhC",
				"Failed to complete the request",
				err.Error()))
		}
	}

	if len(insertFile) != 0 {
		body, err = ioutil.ReadFile(insertFile)
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"IKKp3jAJ",
				"Failed to read the 'insert' file",
				err.Error()))
		}
		_, err = db.Exec(string(body))
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s [%s]",
				"mF0TOjBm",
				"Failed to complete the request",
				err.Error()))
		}
	}

	return nil
}
