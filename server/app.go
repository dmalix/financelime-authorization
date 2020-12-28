/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	packageAuthorization "github.com/dmalix/financelime-authorization/packages/authorization"
	packageAuthorizationAPI "github.com/dmalix/financelime-authorization/packages/authorization/api"
	packageAuthorizationAPIMiddleware "github.com/dmalix/financelime-authorization/packages/authorization/api/middleware"
	packageAuthorizationRepository "github.com/dmalix/financelime-authorization/packages/authorization/repository"
	packageAuthorizationService "github.com/dmalix/financelime-authorization/packages/authorization/service"
	packageSystem "github.com/dmalix/financelime-authorization/packages/system"
	packageSystemAPI "github.com/dmalix/financelime-authorization/packages/system/api"
	packageSystemService "github.com/dmalix/financelime-authorization/packages/system/service"
	"github.com/dmalix/financelime-authorization/utils/cryptographer"
	"github.com/dmalix/financelime-authorization/utils/email"
	"github.com/dmalix/financelime-authorization/utils/jwt"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version   = "unset"
	buildTime = "unset"
	commit    = "unset"
	compiler  = "unset"
)

type App struct {
	httpPort                   string
	emailMessageSenderDaemon   email.EmailSenderDaemon
	httpServer                 *http.Server
	authorizationAPIMiddleware packageAuthorization.APIMiddleware
	authorizationService       packageAuthorization.Service
	authorizationAPI           packageAuthorization.API
	systemService              packageSystem.Service
	systemAPI                  packageSystem.API
}

func NewApp() (*App, error) {

	var (
		app               *App
		err               error
		dbAuthMain        *sql.DB
		dbAuthRead        *sql.DB
		dbBlade           *sql.DB
		config            cfg
		languageContent   models.LanguageContent
		emailMessageQueue = make(chan models.EmailMessage, 300)
	)

	// Init the Config and the Language Content
	// ----------------------------------------

	config, err = initConfig()
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				trace.GetCurrentPoint(),
				"Configuration initialization error",
				err.Error()))
	}

	languageContent, err = initLanguageContent()
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s %s [%s]",
				trace.GetCurrentPoint(),
				"Error initializing language content",
				err.Error()))
	}

	// Databases
	// ---------

	dbAuthMain, err = packageAuthorizationRepository.NewPostgreDB(packageAuthorizationRepository.Config{
		Host:     config.db.authMain.connect.host,
		Port:     config.db.authMain.connect.port,
		SSLMode:  config.db.authMain.connect.sslMode,
		DBName:   config.db.authMain.connect.dbName,
		User:     config.db.authMain.connect.user,
		Password: config.db.authMain.connect.password,
	})
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				trace.GetCurrentPoint(),
				"Failed to init the AuthMain DB",
				err.Error()))
	}

	dbAuthRead, err = packageAuthorizationRepository.NewPostgreDB(packageAuthorizationRepository.Config{
		Host:     config.db.authRead.connect.host,
		Port:     config.db.authRead.connect.port,
		SSLMode:  config.db.authRead.connect.sslMode,
		DBName:   config.db.authRead.connect.dbName,
		User:     config.db.authRead.connect.user,
		Password: config.db.authRead.connect.password,
	})
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				trace.GetCurrentPoint(),
				"Failed to init the AuthRead DB",
				err.Error()))
	}

	dbBlade, err = packageAuthorizationRepository.NewPostgreDB(packageAuthorizationRepository.Config{
		Host:     config.db.blade.connect.host,
		Port:     config.db.blade.connect.port,
		SSLMode:  config.db.blade.connect.sslMode,
		DBName:   config.db.blade.connect.dbName,
		User:     config.db.blade.connect.user,
		Password: config.db.blade.connect.password,
	})
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				trace.GetCurrentPoint(),
				"Failed to init the Blade DB",
				err.Error()))
	}

	err = migrate(dbAuthMain,
		config.db.authMain.migrate.dropFile,
		config.db.authMain.migrate.createFile,
		config.db.authMain.migrate.insertFile)
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				trace.GetCurrentPoint(),
				"Failed to migrate the AuthMain DB",
				err.Error()))
	}

	err = migrate(dbBlade,
		config.db.blade.migrate.dropFile,
		config.db.blade.migrate.createFile,
		"")
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				trace.GetCurrentPoint(),
				"Failed to migrate the Blade DB",
				err.Error()))
	}

	//    Cryptographer
	// --------------------

	cryptoManager := cryptographer.NewCryptographer(config.jwt.secretKey)

	//    JWT
	// -----------

	jwtManager := jwt.NewToken(
		config.jwt.secretKey,
		config.jwt.signingAlgorithm,
		config.jwt.issuer,
		config.jwt.subject,
		config.jwt.accessTokenLifetime,
		config.jwt.refreshTokenLifetime)

	// Email Message
	// -------------

	emailMessageSenderDaemon := email.NewSenderDaemon(
		config.smtp.user,
		config.smtp.password,
		config.smtp.host,
		config.smtp.port,
		emailMessageQueue)

	emailMessageManager := email.NewManager(
		config.mailMessage.from)

	// Authorization
	// ------------------------

	authorizationAPIMiddlewareConfig := packageAuthorizationAPIMiddleware.Config{
		RequestIDRequired: true,
		RequestIDCheck:    true,
	}

	authorizationAPIMiddleware := packageAuthorizationAPIMiddleware.NewMiddleware(
		authorizationAPIMiddlewareConfig,
		jwtManager)

	authorizationRepositoryConfig := packageAuthorizationRepository.ConfigRepository{
		CryptoSalt:              config.crypto.salt,
		JwtRefreshTokenLifetime: config.jwt.refreshTokenLifetime,
	}

	authorizationRepository := packageAuthorizationRepository.NewRepository(
		authorizationRepositoryConfig,
		dbAuthMain,
		dbAuthRead,
		dbBlade)

	authorizationServiceConfig := packageAuthorizationService.ConfigService{
		DomainAPP:              config.domain.app,
		DomainAPI:              config.domain.api,
		AuthInviteCodeRequired: config.auth.inviteCodeRequired,
		CryptoSalt:             config.crypto.salt,
	}

	authorizationService := packageAuthorizationService.NewService(
		authorizationServiceConfig,
		languageContent,
		emailMessageQueue,
		emailMessageManager,
		authorizationRepository,
		cryptoManager,
		jwtManager)

	authorizationAPI := packageAuthorizationAPI.NewHandler(
		authorizationService)

	// System
	// --------------

	systemService := packageSystemService.NewService(
		version,
		buildTime,
		commit,
		compiler)

	systemAPI := packageSystemAPI.NewHandler(
		systemService)

	// Implementation of prepared objects into the REST-API application
	// ----------------------------------------------------------------

	app = &App{
		httpPort:                   config.http.port,
		emailMessageSenderDaemon:   emailMessageSenderDaemon,
		authorizationAPIMiddleware: authorizationAPIMiddleware,
		authorizationService:       authorizationService,
		authorizationAPI:           authorizationAPI,
		systemService:              systemService,
		systemAPI:                  systemAPI,
	}

	return app, nil
}

func (app *App) Run() error {

	// Start the Mail-Sender daemon
	// ----------------------------

	go app.emailMessageSenderDaemon.Run()

	// Start the REST-API application
	// ------------------------------

	mux := http.NewServeMux()

	packageAuthorizationAPI.Router(mux, app.authorizationAPI, app.authorizationAPIMiddleware)
	packageSystemAPI.Router(mux, app.systemAPI, app.authorizationAPIMiddleware)

	app.httpServer = &http.Server{
		Addr:           ":" + app.httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("%s: %s %s [%s]", "FATAL", trace.GetCurrentPoint(), "Failed to listen and serve", err)
		}
	}()

	// Graceful shutdown
	// -----------------

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Printf("%s: %s %s", "INFO", trace.GetCurrentPoint(), "Got SIGINT...")
	case syscall.SIGTERM:
		log.Printf("%s: %s %s", "INFO", trace.GetCurrentPoint(), "Got SIGTERM...")
	}
	log.Printf("%s: %s %s", "INFO", trace.GetCurrentPoint(), "The service is shutting down...")
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	err := app.httpServer.Shutdown(ctx)
	if err == nil {
		log.Printf("%s: %s %s", "INFO", trace.GetCurrentPoint(), "Done")
	}

	return err
}
