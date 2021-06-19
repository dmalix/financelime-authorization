package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dmalix/authorization-service/app/authorization"
	authorizationModel "github.com/dmalix/authorization-service/app/authorization/model"
	authorizationRepository "github.com/dmalix/authorization-service/app/authorization/repository"
	authorizationRestV1 "github.com/dmalix/authorization-service/app/authorization/restv1"
	authorizationService "github.com/dmalix/authorization-service/app/authorization/service"
	"github.com/dmalix/authorization-service/app/information"
	informationRest "github.com/dmalix/authorization-service/app/information/rest"
	informationService "github.com/dmalix/authorization-service/app/information/service"
	"github.com/dmalix/authorization-service/config"
	"github.com/dmalix/jwt"
	"github.com/dmalix/middleware"
	"github.com/dmalix/secretdata"
	"github.com/dmalix/sendmail"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type App struct {
	httpPort                 int
	closeDB                  func() error
	emailMessageSenderDaemon sendmail.Daemon
	httpServer               *http.Server
	commonMiddleware         middleware.Middleware
	authRestV1               authorization.RestV1
	authService              authorization.Service
	infoRest                 information.Rest
	infoService              information.Service
}

func NewApp(logger *zap.Logger, version config.Version) (*App, error) {

	var (
		app                *App
		err                error
		dbAuthMain         *sql.DB
		dbAuthRead         *sql.DB
		dbBlade            *sql.DB
		appConfig          config.App
		appLanguageContent config.LanguageContent
		// TODO Move the number of messages in the queue to configs
		emailMessageQueue = make(chan sendmail.MessageBox, 300)
	)

	// Init Config and Language Content
	appConfig, err = config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to init the config: %s", err)
	}
	appLanguageContent, err = config.InitLanguageContent(appConfig.LanguageContent.File)
	if err != nil {
		return nil, fmt.Errorf("failed to init the language content: %s", err)
	}
	logger.Info("Configuration initialized successfully")

	// Databases
	dbAuthMain, err = authorizationRepository.NewPostgresDB(logger, authorizationModel.ConfigPostgresDB{
		Host:     appConfig.Db.AuthMain.Connect.Host,
		Port:     appConfig.Db.AuthMain.Connect.Port,
		SSLMode:  appConfig.Db.AuthMain.Connect.SslMode,
		DBName:   appConfig.Db.AuthMain.Connect.DbName,
		User:     appConfig.Db.AuthMain.Connect.User,
		Password: appConfig.Db.AuthMain.Connect.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init the AuthMain DB: %s", err)
	}
	dbAuthRead, err = authorizationRepository.NewPostgresDB(logger, authorizationModel.ConfigPostgresDB{
		Host:     appConfig.Db.AuthRead.Connect.Host,
		Port:     appConfig.Db.AuthRead.Connect.Port,
		SSLMode:  appConfig.Db.AuthRead.Connect.SslMode,
		DBName:   appConfig.Db.AuthRead.Connect.DbName,
		User:     appConfig.Db.AuthRead.Connect.User,
		Password: appConfig.Db.AuthRead.Connect.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init the AuthRead DB: %s", err)
	}
	dbBlade, err = authorizationRepository.NewPostgresDB(logger, authorizationModel.ConfigPostgresDB{
		Host:     appConfig.Db.Blade.Connect.Host,
		Port:     appConfig.Db.Blade.Connect.Port,
		SSLMode:  appConfig.Db.Blade.Connect.SslMode,
		DBName:   appConfig.Db.Blade.Connect.DbName,
		User:     appConfig.Db.Blade.Connect.User,
		Password: appConfig.Db.Blade.Connect.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init the Blade DB: %s", err)
	}
	err = config.Migrate(dbAuthMain,
		appConfig.Db.AuthMain.Migration.DropFile,
		appConfig.Db.AuthMain.Migration.CreateFile,
		appConfig.Db.AuthMain.Migration.InsertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate the AuthMain DB: %s", err)
	}
	err = config.Migrate(dbBlade,
		appConfig.Db.Blade.Migration.DropFile,
		appConfig.Db.Blade.Migration.CreateFile,
		"")
	if err != nil {
		return nil, fmt.Errorf("failed to migrate the Blade DB: %s", err)
	}
	closeDB := func() error {
		if err := dbAuthMain.Close(); err != nil {
			return fmt.Errorf("failed to close AuthMain DB: %s", err)
		}
		if err := dbAuthRead.Close(); err != nil {
			return fmt.Errorf("failed to close AuthRead DB: %s", err)
		}
		if err := dbBlade.Close(); err != nil {
			return fmt.Errorf("failed to close Blade DB: %s", err)
		}
		return nil
	}
	logger.Info("Databases initialized successfully")

	// Secret Data
	dataAccess := secretdata.NewSecretData(appConfig.Jwt.AccessSecretKey)
	dataRefresh := secretdata.NewSecretData(appConfig.Jwt.RefreshSecretKey)

	// JWT
	jwtAccess, err := jwt.NewToken(jwt.Config{
		Headers: jwt.Headers{
			Type:               jwt.TokenType,
			SignatureAlgorithm: appConfig.Jwt.AccessSignatureAlgorithm,
		},
		Claims: jwt.Claims{
			Issuer:   appConfig.Jwt.Issuer,
			Audience: appConfig.Jwt.AccessAudience,
			Subject:  jwt.TokenUseAccess,
		},
		ParseOptions: jwt.ParseOptions{
			RequiredClaimIssuer:   true,
			RequiredClaimSubject:  true,
			RequiredClaimJwtID:    true,
			RequiredClaimAudience: true,
			RequiredClaimData:     true,
		},
		TokenLifetimeSec: appConfig.Jwt.AccessTokenLifetime,
		Key:              appConfig.Jwt.AccessSecretKey,
	})
	jwtRefresh, err := jwt.NewToken(jwt.Config{
		Headers: jwt.Headers{
			Type:               jwt.TokenType,
			SignatureAlgorithm: appConfig.Jwt.RefreshSignatureAlgorithm,
		},
		Claims: jwt.Claims{
			Issuer:   appConfig.Jwt.Issuer,
			Audience: appConfig.Jwt.RefreshAudience,
			Subject:  jwt.TokenUseRefresh,
		},
		ParseOptions: jwt.ParseOptions{
			RequiredClaimIssuer:   true,
			RequiredClaimSubject:  true,
			RequiredClaimJwtID:    true,
			RequiredClaimAudience: true,
			RequiredClaimData:     true,
		},
		TokenLifetimeSec: appConfig.Jwt.RefreshTokenLifetime,
		Key:              appConfig.Jwt.RefreshSecretKey,
	})

	// Email Message
	sendMailDaemon := sendmail.NewDaemon(
		appConfig.Smtp.User,
		appConfig.Smtp.Password,
		appConfig.Smtp.Host,
		appConfig.Smtp.Port,
		emailMessageQueue)
	sendMailManager := sendmail.NewManager(appConfig.MailMessage.From)

	// Middleware
	middlewareConfig := middleware.ConfigMiddleware{
		RequestIDRequired:   true,
		RequestIDCheck:      true,
		ContainsJWTAudience: "Authorization",
		DecryptJWTData:      appConfig.Jwt.AccessEncryptData,
		DevelopmentMode:     version.DevelopmentMode,
	}
	commonMiddleware := middleware.NewMiddleware(
		middlewareConfig,
		jwtAccess,
		dataAccess)
	contextGetter := middleware.NewContextGetter()

	// Authorization
	authRepoConfig := authorizationModel.ConfigRepository{
		CryptoSalt:              appConfig.Crypto.Salt,
		JwtRefreshTokenLifetime: appConfig.Jwt.RefreshTokenLifetime,
	}
	authRepo := authorizationRepository.NewRepository(
		authRepoConfig,
		contextGetter,
		dbAuthMain,
		dbAuthRead,
		dbBlade)
	authServiceConfig := authorizationModel.ConfigService{
		DomainAPP:              appConfig.Domain.App,
		DomainAPI:              appConfig.Domain.Api,
		AuthInviteCodeRequired: appConfig.Auth.InviteCodeRequired,
		CryptoSalt:             appConfig.Crypto.Salt,
		JwtAccessEncryptData:   appConfig.Jwt.AccessEncryptData,
		JwtRefreshEncryptData:  appConfig.Jwt.RefreshEncryptData,
	}
	authService := authorizationService.NewService(
		authServiceConfig,
		contextGetter,
		appLanguageContent,
		emailMessageQueue,
		sendMailManager,
		authRepo,
		dataAccess,
		dataRefresh,
		jwtAccess,
		jwtRefresh)
	authRestV1Config := authorizationModel.ConfigRest{
		API:       "v1",
		DomainAPI: appConfig.Domain.Api,
	}
	authRestV1 := authorizationRestV1.NewRest(
		authRestV1Config,
		contextGetter,
		authService)

	// Information
	infoService := informationService.NewService(
		version.Number,
		version.BuildTime,
		version.Commit,
		version.Compiler)
	infoRest := informationRest.NewRest(
		contextGetter,
		infoService)

	// Implementation of prepared objects into the application
	app = &App{
		httpPort:                 appConfig.Http.Port,
		closeDB:                  closeDB,
		emailMessageSenderDaemon: sendMailDaemon,
		commonMiddleware:         commonMiddleware,
		authRestV1:               authRestV1,
		authService:              authService,
		infoRest:                 infoRest,
		infoService:              infoService,
	}

	return app, nil
}

func (app *App) Run(ctx context.Context, logger *zap.Logger) error {

	// Start the Mail-Sender daemon

	go app.emailMessageSenderDaemon.Run(ctx, logger.Named("emailSender"))

	// Start application

	router := mux.NewRouter()
	router.Use(app.commonMiddleware.RemoteAddr(logger.Named("middlewareRemoteAddr")))
	router.Use(app.commonMiddleware.RequestID(logger.Named("middlewareRequestID")))
	router.Use(app.commonMiddleware.Logging(logger.Named("middlewareLogging")))

	routerInfo := router.PathPrefix("/info").Subrouter()
	informationRest.Router(logger.Named("RouterInformation"), routerInfo, app.infoRest)

	routerV1 := router.PathPrefix("/v1").Subrouter()
	authorizationRestV1.RouterV1(logger.Named("RouterV1Authorization"), routerV1, app.authRestV1, app.commonMiddleware)

	app.httpServer = &http.Server{
		Addr:           ":" + strconv.Itoa(app.httpPort),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			switch err {
			case http.ErrServerClosed:
				logger.Info(err.Error())
			default:
				logger.Fatal("Failed to listen and serve", zap.Error(err))
			}
		}
	}()

	logger.Info("Service started")

	// Graceful shutdown

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		logger.Info("Got SIGINT...")
	case syscall.SIGKILL:
		logger.Info("Got SIGKILL...")
	case syscall.SIGQUIT:
		logger.Info("Got SIGQUIT...")
	case syscall.SIGTERM:
		logger.Info("Got SIGTERM...")
	default:
		logger.Info("Undefined killSignal...")
	}
	logger.Info("HTTP service is shutting down...")

	ctxHttpServer, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()
	if err := app.httpServer.Shutdown(ctxHttpServer); err != nil {
		logger.Error("Failed shutdown of the HTTP server", zap.Error(err))
		return err
	}
	logger.Info("HTTP Server is off")

	logger.Info("Databases are closing...")
	if err := app.closeDB(); err != nil {
		logger.Info("Failed to close the databases", zap.Error(err))
	}
	logger.Info("Databases closed")

	return nil
}
