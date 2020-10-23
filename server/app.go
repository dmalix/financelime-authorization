/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	authorizationAPI "github.com/dmalix/financelime-rest-api/packages/authorization/api"
	authorizationDomain "github.com/dmalix/financelime-rest-api/packages/authorization/domain"
	authorizationRepo "github.com/dmalix/financelime-rest-api/packages/authorization/repo"
	authorizationService "github.com/dmalix/financelime-rest-api/packages/authorization/service"
	serverConfig "github.com/dmalix/financelime-rest-api/server/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config      serverConfig.Cfg
	httpServer  *http.Server
	authService authorizationDomain.AccountService
}

func NewApp() (*App, error) {

	var (
		app        *App
		err        error
		dbAuthMain *sql.DB
		dbAuthRead *sql.DB
		dbBlade    *sql.DB
		config     serverConfig.Cfg
	)

	config, err = serverConfig.Init()
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"1oC8sdm0",
				"Configuration initialization error",
				err.Error()))
	}

	dbAuthMain, err = authorizationRepo.NewPostgreDB(authorizationRepo.Config{
		Host:     config.DB.AuthMain.Connect.Host,
		Port:     config.DB.AuthMain.Connect.Port,
		SSLMode:  config.DB.AuthMain.Connect.SSLMode,
		DBName:   config.DB.AuthMain.Connect.DBName,
		User:     config.DB.AuthMain.Connect.User,
		Password: config.DB.AuthMain.Connect.Password,
	})
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"Xp0kMb0D",
				"Failed to init the AuthMain DB",
				err.Error()))
	}

	dbAuthRead, err = authorizationRepo.NewPostgreDB(authorizationRepo.Config{
		Host:     config.DB.AuthRead.Connect.Host,
		Port:     config.DB.AuthRead.Connect.Port,
		SSLMode:  config.DB.AuthRead.Connect.SSLMode,
		DBName:   config.DB.AuthRead.Connect.DBName,
		User:     config.DB.AuthRead.Connect.User,
		Password: config.DB.AuthRead.Connect.Password,
	})
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"tqNC8euP",
				"Failed to init the AuthRead DB",
				err.Error()))
	}

	dbBlade, err = authorizationRepo.NewPostgreDB(authorizationRepo.Config{
		Host:     config.DB.Blade.Connect.Host,
		Port:     config.DB.Blade.Connect.Port,
		SSLMode:  config.DB.Blade.Connect.SSLMode,
		DBName:   config.DB.Blade.Connect.DBName,
		User:     config.DB.Blade.Connect.User,
		Password: config.DB.Blade.Connect.Password,
	})
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"kACb3m8t",
				"Failed to init the Blade DB",
				err.Error()))
	}

	err = serverConfig.Migrate(dbAuthMain,
		config.DB.AuthMain.Migrate.DropFile,
		config.DB.AuthMain.Migrate.CreateFile,
		config.DB.AuthMain.Migrate.InsertFile)
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"sd1oC8m0",
				"Failed to migrate the AuthMain DB",
				err.Error()))
	}

	err = serverConfig.Migrate(dbBlade,
		config.DB.Blade.Migrate.DropFile,
		config.DB.Blade.Migrate.CreateFile,
		"")
	if err != nil {
		return app,
			errors.New(fmt.Sprintf("%s: %s [%s]",
				"uzFsIl4o",
				"Failed to migrate the Blade DB",
				err.Error()))
	}

	accountRepo := authorizationRepo.NewAuthorizationRepo(dbAuthMain, dbAuthRead, dbBlade)

	return &App{
		config: config,
		authService: authorizationService.NewAuthorizationService(
			accountRepo,
			config.Auth.InviteCodeRequired,
		),
	}, nil
}

func (a *App) Run() error {

	mux := http.NewServeMux()

	authorizationAPI.AddRoutes(mux, a.authService)

	a.httpServer = &http.Server{
		Addr:           ":" + a.config.Http.Port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
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
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		log.Print("Done")
	}

	return err
}
