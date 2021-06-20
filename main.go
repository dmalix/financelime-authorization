package main

import (
	"context"
	authorizationApp "github.com/dmalix/authorization-service/app"
	"github.com/dmalix/authorization-service/config"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// @title Authorization Service
// @version v0.3.0-beta
// @description Authorization Service RESTful API service
// @contact.name API Support
// @contact.email dmalix@financelime.com
// @license.name MIT
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

	if version.DevelopmentMode {
		loggerConfig.Development = true
		loggerConfig.Level.SetLevel(zap.DebugLevel)
	} else {
		loggerConfig.Level.SetLevel(zap.InfoLevel)
	}

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

	logger.Info("Welcome to the Authorization Service", zap.Any("version", version))

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
