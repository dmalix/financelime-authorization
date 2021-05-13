/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package main

import (
	"context"
	authorizationApp "github.com/dmalix/financelime-authorization/app"
	"github.com/dmalix/financelime-authorization/config"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// @title Financelime Authorization
// @version v0.2.0-alpha
// @description Financelime Authorization RESTful API service
// @contact.name API Support
// @contact.email dmalix@financelime.com
// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @host api.auth.dev.financelime.com
// @securityDefinitions.apikey authorization
// @in header
// @name authorization
// @schemes https
// @BasePath /

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	version, err := config.InitVersion()
	if err != nil {
		log.Fatalln("failed to build logger", err)
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Development = version.DevelopmentMode

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalln("failed to build logger", err)
	}
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			log.Fatalln("failed to sync logger", err)
		}
	}(logger)

	logger = logger.Named("main")

	logger.Info("Welcome to the Financelime authorization service", zap.Any("version", version))

	loggerApp := logger.Named("app")

	app, err := authorizationApp.NewApp(logger, version)
	if err != nil {
		logger.Fatal("failed to get a new App", zap.Error(err))
	}
	logger.Info("Service successfully initialized")

	err = app.Run(ctx, loggerApp)
	if err != nil {
		logger.Fatal("failed to run the App", zap.Error(err))
	}
	logger.Info("Service stopped")

}
