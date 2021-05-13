/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/authorization"
	"github.com/dmalix/financelime-authorization/packages/cryptographer"
	"github.com/dmalix/financelime-authorization/packages/email"
	"github.com/dmalix/financelime-authorization/packages/jwt"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"github.com/dmalix/financelime-authorization/system"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	versionNumber    = "unset"
	versionBuildTime = "unset"
	versionCommit    = "unset"
	versionCompiler  = "unset"
)

type App struct {
	httpPort                 string
	emailMessageSenderDaemon email.EmailSenderDaemon
	httpServer               *http.Server
	authAPI                  authorization.API
	authAPIMiddleware        middleware.APIMiddleware
	authService              authorization.Service
	sysAPI                   system.API
	sysService               system.Service
}

func NewApp() (*App, error) {

	var (
		app             *App
		err             error
		dbAuthMain      *sql.DB
		dbAuthRead      *sql.DB
		dbBlade         *sql.DB
		config          cfg
		languageContent authorization.LanguageContent
		// TODO Move the number of messages in the queue to configs
		emailMessageQueue = make(chan email.EmailMessage, 300)
	)

	/*****************************************************\
	|           Init Config and Language Content          |
	\*****************************************************/

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

	/*****************************************************\
	|                     Databases                       |
	\*****************************************************/

	dbAuthMain, err = authorization.NewPostgreDB(authorization.ConfigPostgreDB{
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

	dbAuthRead, err = authorization.NewPostgreDB(authorization.ConfigPostgreDB{
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

	dbBlade, err = authorization.NewPostgreDB(authorization.ConfigPostgreDB{
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

	/*****************************************************\
	|                    Cryptographer                    |
	\*****************************************************/

	cryptoManager := cryptographer.NewCryptographer(config.jwt.secretKey)

	/*****************************************************\
	|                         JWT                         |
	\*****************************************************/

	jwtManager := jwt.NewToken(
		config.jwt.secretKey,
		config.jwt.signingAlgorithm,
		config.jwt.issuer,
		config.jwt.subject,
		config.jwt.accessTokenLifetime,
		config.jwt.refreshTokenLifetime)

	/*****************************************************\
	|                     Email Message                   |
	\*****************************************************/

	emailMessageSenderDaemon := email.NewSenderDaemon(
		config.smtp.user,
		config.smtp.password,
		config.smtp.host,
		config.smtp.port,
		emailMessageQueue)

	emailMessageManager := email.NewManager(
		config.mailMessage.from)

	/*****************************************************\
	|                     Middleware                      |
	\*****************************************************/

	middlewareConfig := middleware.ConfigMiddleware{
		RequestIDRequired: true,
		RequestIDCheck:    true,
	}

	commonMiddleware := middleware.NewMiddleware(
		middlewareConfig,
		jwtManager)

	/*****************************************************\
	|                     Authorization                   |
	\*****************************************************/

	authRepoConfig := authorization.ConfigRepository{
		CryptoSalt:              config.crypto.salt,
		JwtRefreshTokenLifetime: config.jwt.refreshTokenLifetime,
	}

	authRepo := authorization.NewRepository(
		authRepoConfig,
		dbAuthMain,
		dbAuthRead,
		dbBlade)

	authorizationServiceConfig := authorization.ConfigService{
		DomainAPP:              config.domain.app,
		DomainAPI:              config.domain.api,
		AuthInviteCodeRequired: config.auth.inviteCodeRequired,
		CryptoSalt:             config.crypto.salt,
	}

	authService := authorization.NewService(
		authorizationServiceConfig,
		languageContent,
		emailMessageQueue,
		emailMessageManager,
		authRepo,
		cryptoManager,
		jwtManager)

	authAPI := authorization.NewAPI(
		authService)

	/*****************************************************\
	|                       System                        |
	\*****************************************************/

	systemService := system.NewService(
		versionNumber,
		versionBuildTime,
		versionCommit,
		versionCompiler)

	systemAPI := system.NewAPI(
		systemService)

	/*****************************************************\
	|          Implementation of prepared objects         |
	|                 into the application                |
	\*****************************************************/

	app = &App{
		httpPort:                 config.http.port,
		emailMessageSenderDaemon: emailMessageSenderDaemon,
		authAPI:                  authAPI,
		authAPIMiddleware:        commonMiddleware,
		authService:              authService,
		sysAPI:                   systemAPI,
		sysService:               systemService,
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) error {

	/*****************************************************\
	|             Start the Mail-Sender daemon            |
	\*****************************************************/

	go app.emailMessageSenderDaemon.Run(ctx)

	/*****************************************************\
	|                  Start application                  |
	\*****************************************************/

	router := mux.NewRouter()
	router.Use(app.authAPIMiddleware.Logging())
	routerV1 := router.PathPrefix("/v1").Subrouter()

	authorization.Router(ctx, routerV1, app.authAPI, app.authAPIMiddleware)
	system.Router(ctx, routerV1, app.sysAPI, app.authAPIMiddleware)

	app.httpServer = &http.Server{
		Addr:           ":" + app.httpPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			switch err {
			case http.ErrServerClosed:
				log.Printf("%s: %s %s", "INFO", trace.GetCurrentPoint(), err.Error())
			default:
				log.Fatalf("%s: %s %s [%s]", "FATAL", trace.GetCurrentPoint(), "failed to listen and serve", err)
			}
		}
	}()

	/*****************************************************\
	|                  Graceful shutdown                  |
	\*****************************************************/

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
	ctxHttpServer, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()
	err := app.httpServer.Shutdown(ctxHttpServer)
	if err != nil {
		log.Println("failed shutdown of the HTTP Server", err)
	}
	log.Printf("%s: %s %s", "INFO", trace.GetCurrentPoint(), "Done")

	return err
}
