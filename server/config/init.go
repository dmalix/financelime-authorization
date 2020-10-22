/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Cfg struct {
	Http struct {
		Port string
	}
	Auth struct {
		InviteCodeRequired bool
	}
	DB struct {
		AuthMain CfgDB
		AuthRead CfgDB
		Blade    CfgDB
	}
}

type CfgDB struct {
	Connect struct {
		Host     string
		Port     string
		SSLMode  string
		DBName   string
		User     string
		Password string
	}
	Migrate struct {
		DropFile   string
		CreateFile string
		InsertFile string
	}
}

func Init() (Cfg, error) {

	var err error
	var config Cfg

	config.Http.Port = os.Getenv("HTTP_SERVER_PORT")
	if len(config.Http.Port) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"o8sdm1C0",
			"The HTTP_SERVER_PORT environment variable has empty value"))
	}

	config.Auth.InviteCodeRequired, err = strconv.ParseBool(os.Getenv("AUTH_INVITE_CODE_REQUIRED"))
	if err != nil {
		return config, errors.New(fmt.Sprintf("%s: %s [%s]",
			"8sdm1oC0",
			"The AUTH_INVITE_CODE_REQUIRED environment variable has no boolean value",
			err.Error()))
	}

	config.DB.AuthMain.Connect.Host = os.Getenv("DB_AUTH_MAIN_CONNECT_HOST")
	config.DB.AuthMain.Connect.Port = os.Getenv("DB_AUTH_MAIN_CONNECT_PORT")
	config.DB.AuthMain.Connect.SSLMode = os.Getenv("DB_AUTH_MAIN_CONNECT_SSLMODE")
	config.DB.AuthMain.Connect.DBName = os.Getenv("DB_AUTH_MAIN_CONNECT_DBNAME")
	config.DB.AuthMain.Connect.User = os.Getenv("DB_AUTH_MAIN_CONNECT_USER")
	config.DB.AuthMain.Connect.Password = os.Getenv("DB_AUTH_MAIN_CONNECT_PASSWORD")
	if len(config.DB.AuthMain.Connect.Host) == 0 || len(config.DB.AuthMain.Connect.Port) == 0 ||
		len(config.DB.AuthMain.Connect.SSLMode) == 0 || len(config.DB.AuthMain.Connect.DBName) == 0 ||
		len(config.DB.AuthMain.Connect.User) == 0 || len(config.DB.AuthMain.Connect.Password) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"HY3kabzX",
			"One or more environment variables are null [DB_AUTH_MAIN_CONNECT_*]"))
	}

	config.DB.AuthRead.Connect.Host = os.Getenv("DB_AUTH_READ_CONNECT_HOST")
	config.DB.AuthRead.Connect.Port = os.Getenv("DB_AUTH_READ_CONNECT_PORT")
	config.DB.AuthRead.Connect.SSLMode = os.Getenv("DB_AUTH_READ_CONNECT_SSLMODE")
	config.DB.AuthRead.Connect.DBName = os.Getenv("DB_AUTH_READ_CONNECT_DBNAME")
	config.DB.AuthRead.Connect.User = os.Getenv("DB_AUTH_READ_CONNECT_USER")
	config.DB.AuthRead.Connect.Password = os.Getenv("DB_AUTH_READ_CONNECT_PASSWORD")
	if len(config.DB.AuthRead.Connect.Host) == 0 || len(config.DB.AuthRead.Connect.Port) == 0 ||
		len(config.DB.AuthRead.Connect.SSLMode) == 0 || len(config.DB.AuthRead.Connect.DBName) == 0 ||
		len(config.DB.AuthRead.Connect.User) == 0 || len(config.DB.AuthRead.Connect.Password) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"Y3kabHzX",
			"One or more environment variables are null [DB_AUTH_READ_CONNECT_*]"))
	}

	config.DB.Blade.Connect.Host = os.Getenv("DB_BLADE_CONNECT_HOST")
	config.DB.Blade.Connect.Port = os.Getenv("DB_BLADE_CONNECT_PORT")
	config.DB.Blade.Connect.SSLMode = os.Getenv("DB_BLADE_CONNECT_SSLMODE")
	config.DB.Blade.Connect.DBName = os.Getenv("DB_BLADE_CONNECT_DBNAME")
	config.DB.Blade.Connect.User = os.Getenv("DB_BLADE_CONNECT_USER")
	config.DB.Blade.Connect.Password = os.Getenv("DB_BLADE_CONNECT_PASSWORD")
	if len(config.DB.Blade.Connect.Host) == 0 || len(config.DB.Blade.Connect.Port) == 0 ||
		len(config.DB.Blade.Connect.SSLMode) == 0 || len(config.DB.Blade.Connect.DBName) == 0 ||
		len(config.DB.Blade.Connect.User) == 0 || len(config.DB.Blade.Connect.Password) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"zY3kabHX",
			"One or more environment variables are null [DB_AUTH_BLADE_CONNECT_*]"))
	}

	config.DB.AuthMain.Migrate.DropFile = os.Getenv("DB_AUTH_MAIN_MIGRATE_DROPFILE")
	config.DB.AuthMain.Migrate.CreateFile = os.Getenv("DB_AUTH_MAIN_MIGRATE_CREATEFILE")
	config.DB.AuthMain.Migrate.InsertFile = os.Getenv("DB_AUTH_MAIN_MIGRATE_INSERTFILE")

	config.DB.Blade.Migrate.DropFile = os.Getenv("DB_BLADE_MIGRATE_DROPFILE")
	config.DB.Blade.Migrate.CreateFile = os.Getenv("DB_BLADE_MIGRATE_CREATEFILE")

	return config, nil
}
