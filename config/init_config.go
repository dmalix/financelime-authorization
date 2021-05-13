/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package config

import (
	"fmt"
	"os"
	"strconv"
)

func InitConfig() (App, error) {

	var err error
	var config App
	const messageEnvironmentVariableIsEmpty = "the %s environment variable is empty"
	const messageEnvironmentVariableIsNotNumber = "the %s environment variable is not a number: %s"
	const messageEnvironmentVariableIsNull = "the %s environment variable is null"
	const messageEnvironmentVariableIsNotBoolean = "the %s environment variable is not boolean: %s"

	// Language content
	if config.LanguageContent.File = os.Getenv(envLanguageContentFile); config.LanguageContent.File == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envLanguageContentFile)
	}

	// Domain
	if config.Domain.App = os.Getenv(envDomainApp); config.Domain.App == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDomainApp)
	}
	if config.Domain.Api = os.Getenv(envDomainApi); config.Domain.Api == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDomainApi)
	}

	// Port
	config.Http.Port, err = strconv.Atoi(os.Getenv(envHttpServerPort))
	if err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotNumber, envHttpServerPort, err)
	}
	if config.Http.Port == 0 {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNull, envHttpServerPort)
	}

	// Invite Code
	if config.Auth.InviteCodeRequired, err = strconv.ParseBool(os.Getenv(envAuthInviteCodeRequired)); err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotBoolean, envAuthInviteCodeRequired, err)
	}

	// DB Auth
	if config.Db.AuthMain.Connect.Host = os.Getenv(envDbAuthMainConnectHost); config.Db.AuthMain.Connect.Host == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainConnectHost)
	}
	config.Db.AuthMain.Connect.Port, err = strconv.Atoi(os.Getenv(envDbAuthMainConnectPort))
	if err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotNumber, envDbAuthMainConnectPort, err)
	}
	if config.Db.AuthMain.Connect.Port == 0 {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNull, envDbAuthMainConnectPort)
	}
	if config.Db.AuthMain.Connect.SslMode = os.Getenv(envDbAuthMainConnectSslMode); config.Db.AuthMain.Connect.SslMode == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainConnectSslMode)
	}
	if config.Db.AuthMain.Connect.DbName = os.Getenv(envDbAuthMainConnectDbName); config.Db.AuthMain.Connect.DbName == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainConnectDbName)
	}
	if config.Db.AuthMain.Connect.User = os.Getenv(envDbAuthMainConnectUser); config.Db.AuthMain.Connect.User == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainConnectUser)
	}
	if config.Db.AuthMain.Connect.Password = os.Getenv(envDbAuthMainConnectPassword); config.Db.AuthMain.Connect.Password == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainConnectPassword)
	}

	// DB AuthRead
	if config.Db.AuthRead.Connect.Host = os.Getenv(envDbAuthReadConnectHost); config.Db.AuthRead.Connect.Host == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthReadConnectHost)
	}
	if config.Db.AuthRead.Connect.Port, err = strconv.Atoi(os.Getenv(envDbAuthReadConnectPort)); err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotNumber, envDbAuthReadConnectPort, err)
	}
	if config.Db.AuthRead.Connect.Port == 0 {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNull, envDbAuthReadConnectPort)
	}
	if config.Db.AuthRead.Connect.SslMode = os.Getenv(envDbAuthReadConnectSslMode); config.Db.AuthRead.Connect.SslMode == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthReadConnectSslMode)
	}
	if config.Db.AuthRead.Connect.DbName = os.Getenv(envDbAuthReadConnectDbName); config.Db.AuthRead.Connect.DbName == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthReadConnectDbName)
	}
	if config.Db.AuthRead.Connect.User = os.Getenv(envDbAuthReadConnectUser); config.Db.AuthRead.Connect.User == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthReadConnectUser)
	}
	if config.Db.AuthRead.Connect.Password = os.Getenv(envDbAuthReadConnectPassword); config.Db.AuthRead.Connect.Password == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthReadConnectPassword)
	}

	// DB AuthMain Migrate
	if config.Db.AuthMain.Migration.DropFile = os.Getenv(envDbAuthMainMigrateDropfile); config.Db.AuthMain.Migration.DropFile == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainMigrateDropfile)
	}
	if config.Db.AuthMain.Migration.CreateFile = os.Getenv(envDbAuthMainMigrateCreateFile); config.Db.AuthMain.Migration.CreateFile == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainMigrateCreateFile)
	}
	if config.Db.AuthMain.Migration.InsertFile = os.Getenv(envDbAuthMainMigrateInsertFile); config.Db.AuthMain.Migration.InsertFile == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbAuthMainMigrateInsertFile)
	}

	// DB Blade
	if config.Db.Blade.Connect.Host = os.Getenv(envDbBladeConnectHost); config.Db.Blade.Connect.Host == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeConnectHost)
	}
	if config.Db.Blade.Connect.Port, err = strconv.Atoi(os.Getenv(envDbBladeConnectPort)); err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotNumber, envDbBladeConnectPort, err)
	}
	if config.Db.Blade.Connect.Port == 0 {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNull, envDbBladeConnectPort)
	}
	if config.Db.Blade.Connect.SslMode = os.Getenv(envDbBladeConnectSslMode); config.Db.Blade.Connect.SslMode == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeConnectSslMode)
	}
	if config.Db.Blade.Connect.DbName = os.Getenv(envDbBladeConnectDbName); config.Db.Blade.Connect.DbName == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeConnectDbName)
	}
	if config.Db.Blade.Connect.User = os.Getenv(envDbBladeConnectUser); config.Db.Blade.Connect.User == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeConnectUser)
	}
	if config.Db.Blade.Connect.Password = os.Getenv(envDbBladeConnectPassword); config.Db.Blade.Connect.Password == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeConnectPassword)
	}

	// DB Blade Migrate
	if config.Db.Blade.Migration.DropFile = os.Getenv(envDbBladeMigrateDropFile); config.Db.Blade.Migration.DropFile == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeMigrateDropFile)
	}
	if config.Db.Blade.Migration.CreateFile = os.Getenv(envDbBladeMigrateCreateFile); config.Db.Blade.Migration.CreateFile == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envDbBladeMigrateCreateFile)
	}

	// SMTP
	if config.Smtp.User = os.Getenv(envSmtpUser); config.Smtp.User == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envSmtpUser)
	}
	if config.Smtp.Password = os.Getenv(envSmtpPassword); config.Smtp.Password == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envSmtpPassword)
	}
	if config.Smtp.Host = os.Getenv(envSmtpHost); config.Smtp.Host == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envSmtpHost)
	}
	if config.Smtp.Port, err = strconv.Atoi(os.Getenv(envSmtpPort)); err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envSmtpPort)
	}
	if config.Smtp.Port == 0 {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNull, envSmtpPort)
	}

	// Mail Message
	if config.MailMessage.From.Name = os.Getenv(envMailMessageFromName); config.MailMessage.From.Name == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envMailMessageFromName)
	}
	if config.MailMessage.From.Address = os.Getenv(envMailMessageFromAddr); config.MailMessage.From.Address == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envMailMessageFromAddr)
	}

	// Crypro
	if config.Crypto.Salt = os.Getenv(envCryptoSalt); config.Crypto.Salt == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envCryptoSalt)
	}

	// JWT
	if config.Jwt.SecretKey = os.Getenv(envJwtSecretKey); config.Jwt.SecretKey == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envJwtSecretKey)
	}
	if config.Jwt.SigningAlgorithm = os.Getenv(envJwtSigningAlgorithm); config.Jwt.SigningAlgorithm == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envJwtSigningAlgorithm)
	}
	if config.Jwt.Issuer = os.Getenv(envJwtIssuer); config.Jwt.Issuer == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envJwtIssuer)
	}
	if config.Jwt.Subject = os.Getenv(envJwtSubject); config.Jwt.Subject == "" {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsEmpty, envJwtSubject)
	}
	if config.Jwt.AccessTokenLifetime, err = strconv.Atoi(os.Getenv(envJwtAccessTokenLifeTime)); err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotNumber, envJwtAccessTokenLifeTime, err)
	}
	if config.Jwt.RefreshTokenLifetime, err = strconv.Atoi(os.Getenv(envJwtRefreshTokenLifeTime)); err != nil {
		return App{}, fmt.Errorf(messageEnvironmentVariableIsNotNumber, envJwtRefreshTokenLifeTime, err)
	}

	return config, nil
}
