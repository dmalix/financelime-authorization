/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type Config struct {
	Host     string
	Port     string
	SSLMode  string
	DBName   string
	User     string
	Password string
}

func NewPostgreDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s user=%s password=%s",
			cfg.Host, cfg.Port, cfg.SSLMode, cfg.DBName, cfg.User, cfg.Password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s [%s]",
			"hAjgNlS8",
			"Failed to open DB connection",
			err.Error()))
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s [%s]",
			"oA0aQs3K",
			"Failed to ping to DB",
			err.Error()))
	}

	return db, nil
}
