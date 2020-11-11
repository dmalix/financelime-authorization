/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package server

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strconv"
)

type cfg struct {
	domain struct {
		app string
		api string
	}
	http struct {
		port string
	}
	auth struct {
		inviteCodeRequired bool
	}
	db struct {
		authMain cfgDB
		authRead cfgDB
		blade    cfgDB
	}
	smtp struct {
		user     string
		password string
		host     string
		port     string
	}
	mailMessage struct {
		from mail.Address
	}
	crypto struct {
		salt string
	}
	jwt struct {
		secretKey            string
		signingAlgorithm     string
		issuer               string
		subject              string
		accessTokenLifetime  int
		refreshTokenLifetime int
	}
}

type cfgDB struct {
	connect struct {
		host     string
		port     string
		sslMode  string
		dbName   string
		user     string
		password string
	}
	migrate struct {
		dropFile   string
		createFile string
		insertFile string
	}
}

func initConfig() (cfg, error) {

	var err error
	var config cfg

	// DOMAIN
	// ------

	config.domain.app = os.Getenv("DOMAIN_APP")
	if len(config.domain.app) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"sd1Co8m0",
			"The DOMAIN_APP environment variable has empty value"))
	}

	config.domain.api = os.Getenv("DOMAIN_API")
	if len(config.domain.api) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"sd1Co8m0",
			"The DOMAIN_API environment variable has empty value"))
	}

	// HTTP
	// ----

	config.http.port = os.Getenv("HTTP_SERVER_PORT")
	if len(config.http.port) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"o8sdm1C0",
			"The HTTP_SERVER_PORT environment variable has empty value"))
	}

	// AUTH
	// ----

	config.auth.inviteCodeRequired, err = strconv.ParseBool(os.Getenv("AUTH_INVITE_CODE_REQUIRED"))
	if err != nil {
		return config, errors.New(fmt.Sprintf("%s: %s [%s]",
			"8sdm1oC0",
			"The AUTH_INVITE_CODE_REQUIRED environment variable has no boolean value",
			err.Error()))
	}

	// DB AUTH
	// -------

	config.db.authMain.connect.host = os.Getenv("DB_AUTH_MAIN_CONNECT_HOST")
	config.db.authMain.connect.port = os.Getenv("DB_AUTH_MAIN_CONNECT_PORT")
	config.db.authMain.connect.sslMode = os.Getenv("DB_AUTH_MAIN_CONNECT_SSLMODE")
	config.db.authMain.connect.dbName = os.Getenv("DB_AUTH_MAIN_CONNECT_DBNAME")
	config.db.authMain.connect.user = os.Getenv("DB_AUTH_MAIN_CONNECT_USER")
	config.db.authMain.connect.password = os.Getenv("DB_AUTH_MAIN_CONNECT_PASSWORD")
	if len(config.db.authMain.connect.host) == 0 || len(config.db.authMain.connect.port) == 0 ||
		len(config.db.authMain.connect.sslMode) == 0 || len(config.db.authMain.connect.dbName) == 0 ||
		len(config.db.authMain.connect.user) == 0 || len(config.db.authMain.connect.password) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"HY3kabzX",
			"One or more environment variables are null [DB_AUTH_MAIN_CONNECT_*]"))
	}

	config.db.authRead.connect.host = os.Getenv("DB_AUTH_READ_CONNECT_HOST")
	config.db.authRead.connect.port = os.Getenv("DB_AUTH_READ_CONNECT_PORT")
	config.db.authRead.connect.sslMode = os.Getenv("DB_AUTH_READ_CONNECT_SSLMODE")
	config.db.authRead.connect.dbName = os.Getenv("DB_AUTH_READ_CONNECT_DBNAME")
	config.db.authRead.connect.user = os.Getenv("DB_AUTH_READ_CONNECT_USER")
	config.db.authRead.connect.password = os.Getenv("DB_AUTH_READ_CONNECT_PASSWORD")
	if len(config.db.authRead.connect.host) == 0 || len(config.db.authRead.connect.port) == 0 ||
		len(config.db.authRead.connect.sslMode) == 0 || len(config.db.authRead.connect.dbName) == 0 ||
		len(config.db.authRead.connect.user) == 0 || len(config.db.authRead.connect.password) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"Y3kabHzX",
			"One or more environment variables are null [DB_AUTH_READ_CONNECT_*]"))
	}

	config.db.authMain.migrate.dropFile = os.Getenv("DB_AUTH_MAIN_MIGRATE_DROPFILE")
	config.db.authMain.migrate.createFile = os.Getenv("DB_AUTH_MAIN_MIGRATE_CREATEFILE")
	config.db.authMain.migrate.insertFile = os.Getenv("DB_AUTH_MAIN_MIGRATE_INSERTFILE")

	// DB BLADE
	// --------

	config.db.blade.connect.host = os.Getenv("DB_BLADE_CONNECT_HOST")
	config.db.blade.connect.port = os.Getenv("DB_BLADE_CONNECT_PORT")
	config.db.blade.connect.sslMode = os.Getenv("DB_BLADE_CONNECT_SSLMODE")
	config.db.blade.connect.dbName = os.Getenv("DB_BLADE_CONNECT_DBNAME")
	config.db.blade.connect.user = os.Getenv("DB_BLADE_CONNECT_USER")
	config.db.blade.connect.password = os.Getenv("DB_BLADE_CONNECT_PASSWORD")
	if len(config.db.blade.connect.host) == 0 || len(config.db.blade.connect.port) == 0 ||
		len(config.db.blade.connect.sslMode) == 0 || len(config.db.blade.connect.dbName) == 0 ||
		len(config.db.blade.connect.user) == 0 || len(config.db.blade.connect.password) == 0 {
		return config, errors.New(fmt.Sprintf("%s: %s",
			"zY3kabHX",
			"One or more environment variables are null [DB_AUTH_BLADE_CONNECT_*]"))
	}

	config.db.blade.migrate.dropFile = os.Getenv("DB_BLADE_MIGRATE_DROPFILE")
	config.db.blade.migrate.createFile = os.Getenv("DB_BLADE_MIGRATE_CREATEFILE")

	// AUTH SMTP
	// ---------

	config.smtp.user = os.Getenv("SMTP_USER")
	config.smtp.password = os.Getenv("SMTP_PASSWORD")
	config.smtp.host = os.Getenv("SMTP_HOST")
	config.smtp.port = os.Getenv("SMTP_PORT")

	// MAIL MESSAGE
	// ------------

	config.mailMessage.from.Name = os.Getenv("MAIL_MESSAGE_FROM_NAME")
	config.mailMessage.from.Address = os.Getenv("MAIL_MESSAGE_FROM_ADDR")

	//   CRYPTO
	// ------------

	config.crypto.salt = os.Getenv("CRYPTO_SALT")

	//     JWT
	// ------------

	config.jwt.secretKey = os.Getenv("JWT_SECRET_KEY")
	config.jwt.signingAlgorithm = os.Getenv("JWT_SIGNING_ALGORITHM")
	config.jwt.issuer = os.Getenv("JWT_ISSUER")
	config.jwt.subject = os.Getenv("JWT_SUBJECT")
	config.jwt.accessTokenLifetime, err = strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_LIFETIME"))
	if err != nil {
		return config, errors.New(fmt.Sprintf("%s: %s [%s]",
			"1oC8sdm0",
			"The JWT_ACCESS_TOKEN_LIFETIME environment variable has no integer type",
			err.Error()))
	}
	config.jwt.refreshTokenLifetime, err = strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_LIFETIME"))
	if err != nil {
		return config, errors.New(fmt.Sprintf("%s: %s [%s]",
			"s1oC8dm0",
			"The JWT_REFRESH_TOKEN_LIFETIME environment variable has no integer type",
			err.Error()))
	}

	return config, nil
}
