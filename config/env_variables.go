package config

const (

	// DOMAIN
	envDomainApp = "DOMAIN_APP"
	envDomainApi = "DOMAIN_API"
	// HTTP
	envHttpServerPort = "HTTP_SERVER_PORT"
	// AUTH
	envAuthInviteCodeRequired = "AUTH_INVITE_CODE_REQUIRED"
	// DB AUTH MAIN
	envDbAuthMainConnectHost       = "DB_AUTH_MAIN_CONNECT_HOST"
	envDbAuthMainConnectPort       = "DB_AUTH_MAIN_CONNECT_PORT"
	envDbAuthMainConnectSslMode    = "DB_AUTH_MAIN_CONNECT_SSLMODE"
	envDbAuthMainConnectDbName     = "DB_AUTH_MAIN_CONNECT_DBNAME"
	envDbAuthMainConnectUser       = "DB_AUTH_MAIN_CONNECT_USER"
	envDbAuthMainConnectPassword   = "DB_AUTH_MAIN_CONNECT_PASSWORD"
	envDbAuthMainMigrateDropfile   = "DB_AUTH_MAIN_MIGRATE_DROPFILE"
	envDbAuthMainMigrateCreateFile = "DB_AUTH_MAIN_MIGRATE_CREATEFILE"
	envDbAuthMainMigrateInsertFile = "DB_AUTH_MAIN_MIGRATE_INSERTFILE"
	// DB AUTH READ
	envDbAuthReadConnectHost     = "DB_AUTH_READ_CONNECT_HOST"
	envDbAuthReadConnectPort     = "DB_AUTH_READ_CONNECT_PORT"
	envDbAuthReadConnectSslMode  = "DB_AUTH_READ_CONNECT_SSLMODE"
	envDbAuthReadConnectDbName   = "DB_AUTH_READ_CONNECT_DBNAME"
	envDbAuthReadConnectUser     = "DB_AUTH_READ_CONNECT_USER"
	envDbAuthReadConnectPassword = "DB_AUTH_READ_CONNECT_PASSWORD"
	// DB BLADE
	envDbBladeConnectHost       = "DB_BLADE_CONNECT_HOST"
	envDbBladeConnectPort       = "DB_BLADE_CONNECT_PORT"
	envDbBladeConnectSslMode    = "DB_BLADE_CONNECT_SSLMODE"
	envDbBladeConnectDbName     = "DB_BLADE_CONNECT_DBNAME"
	envDbBladeConnectUser       = "DB_BLADE_CONNECT_USER"
	envDbBladeConnectPassword   = "DB_BLADE_CONNECT_PASSWORD"
	envDbBladeMigrateDropFile   = "DB_BLADE_MIGRATE_DROPFILE"
	envDbBladeMigrateCreateFile = "DB_BLADE_MIGRATE_CREATEFILE"
	// SMTP
	envSmtpUser     = "SMTP_USER"
	envSmtpPassword = "SMTP_PASSWORD"
	envSmtpHost     = "SMTP_HOST"
	envSmtpPort     = "SMTP_PORT"
	// MAIL MESSAGE
	envMailMessageFromName = "MAIL_MESSAGE_FROM_NAME"
	envMailMessageFromAddr = "MAIL_MESSAGE_FROM_ADDR"
	// CONTENT
	envLanguageContentFile = "LANGUAGE_CONTENT_FILE"
	// CRYPTO
	envCryptoSalt = "CRYPTO_SALT"
	// JWT
	envJwtIssuer                    = "JWT_ISSUER"
	envJwtAccessAudience            = "JWT_ACCESS_AUDIENCE"
	envJwtAccessSecretKey           = "JWT_ACCESS_SECRET_KEY"
	envJwtAccessSignatureAlgorithm  = "JWT_ACCESS_SIGNATURE_ALGORITHM"
	envJwtAccessEncryptData         = "JWT_ACCESS_ENCRYPT_DATA"
	envJwtAccessTokenLifeTime       = "JWT_ACCESS_TOKEN_LIFETIME"
	envJwtRefreshAudience           = "JWT_REFRESH_AUDIENCE"
	envJwtRefreshSecretKey          = "JWT_REFRESH_SECRET_KEY"
	envJwtRefreshSignatureAlgorithm = "JWT_REFRESH_SIGNATURE_ALGORITHM"
	envJwtRefreshEncryptData        = "JWT_REFRESH_ENCRYPT_DATA"
	envJwtRefreshTokenLifeTime      = "JWT_REFRESH_TOKEN_LIFETIME"
)
