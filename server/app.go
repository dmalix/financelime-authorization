/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/packages/authorization"
	authorizationAPI "github.com/dmalix/financelime-rest-api/packages/authorization/api"
	authorizationMiddleware "github.com/dmalix/financelime-rest-api/packages/authorization/api/middleware"
	authorizationRepo "github.com/dmalix/financelime-rest-api/packages/authorization/repository"
	"github.com/dmalix/financelime-rest-api/packages/authorization/service"
	emailMessage "github.com/dmalix/financelime-rest-api/utils/email"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type emailSenderDaemon interface {
	Run()
}

type App struct {
	httpPort                 string
	emailMessageSenderDaemon emailSenderDaemon
	httpServer               *http.Server
	authService              authorization.Service
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

	config, err = initConfig()
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"1oC8sdm0",
				"Configuration initialization error",
				err.Error()))
	}

	dbAuthMain, err = authorizationRepo.NewPostgreDB(authorizationRepo.Config{
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
				"Xp0kMb0D",
				"Failed to init the AuthMain DB",
				err.Error()))
	}

	dbAuthRead, err = authorizationRepo.NewPostgreDB(authorizationRepo.Config{
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
				"tqNC8euP",
				"Failed to init the AuthRead DB",
				err.Error()))
	}

	dbBlade, err = authorizationRepo.NewPostgreDB(authorizationRepo.Config{
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
				"kACb3m8t",
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
				"sd1oC8m0",
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
				"uzFsIl4o",
				"Failed to migrate the Blade DB",
				err.Error()))
	}

	userRepo := authorizationRepo.NewRepository(dbAuthMain, dbAuthRead, dbBlade)
	emailMessageSenderDaemon := emailMessage.NewSenderDaemon(config.smtp.user, config.smtp.password,
		config.smtp.host, config.smtp.port, emailMessageQueue)
	userMessage := emailMessage.NewManager(config.mailMessage.from)

	languageContent, err = initLanguageContent()
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"d1oC8sm0",
				"Error initializing language content",
				err.Error()))
	}

	return &App{
		httpPort:                 config.http.port,
		emailMessageSenderDaemon: emailMessageSenderDaemon,
		authService: service.NewService(
			config.domain.api,
			config.auth.inviteCodeRequired,
			languageContent,
			emailMessageQueue,
			userMessage,
			userRepo),
	}, nil
}

func (app *App) Run() error {

	// Start mail sender daemon
	// ------------------------

	go app.emailMessageSenderDaemon.Run()

	// Start REST-API
	// --------------

	mux := http.NewServeMux()

	authorizationAPI.Router(mux, app.authService, authorizationMiddleware.RequestID)

	app.httpServer = &http.Server{
		Addr:           ":" + app.httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("FATAL [D49VshMa: Failed to listen and serve: [%v]]", err)
		}
	}()

	// Graceful shutdown

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}
	log.Print("The service is shutting down...")
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	err := app.httpServer.Shutdown(ctx)
	if err != nil {
		log.Print("Done")
	}

	return err
}
