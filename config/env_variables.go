package config

const (
	envLanguageContentFile = "LANGUAGE_CONTENT_FILE"

	envDomainApp = "DOMAIN_APP"
	envDomainApi = "DOMAIN_API"

	envHttpServerPort = "HTTP_SERVER_PORT"

	envAuthInviteCodeRequired = "AUTH_INVITE_CODE_REQUIRED"

	envDbAuthMainConnectHost     = "DB_AUTH_MAIN_CONNECT_HOST"
	envDbAuthMainConnectPort     = "DB_AUTH_MAIN_CONNECT_PORT"
	envDbAuthMainConnectSslMode  = "DB_AUTH_MAIN_CONNECT_SSLMODE"
	envDbAuthMainConnectDbName   = "DB_AUTH_MAIN_CONNECT_DBNAME"
	envDbAuthMainConnectUser     = "DB_AUTH_MAIN_CONNECT_USER"
	envDbAuthMainConnectPassword = "DB_AUTH_MAIN_CONNECT_PASSWORD"

	envDbAuthReadConnectHost     = "DB_AUTH_READ_CONNECT_HOST"
	envDbAuthReadConnectPort     = "DB_AUTH_READ_CONNECT_PORT"
	envDbAuthReadConnectSslMode  = "DB_AUTH_READ_CONNECT_SSLMODE"
	envDbAuthReadConnectDbName   = "DB_AUTH_READ_CONNECT_DBNAME"
	envDbAuthReadConnectUser     = "DB_AUTH_READ_CONNECT_USER"
	envDbAuthReadConnectPassword = "DB_AUTH_READ_CONNECT_PASSWORD"

	envDbAuthMainMigrateDropfile   = "DB_AUTH_MAIN_MIGRATE_DROPFILE"
	envDbAuthMainMigrateCreateFile = "DB_AUTH_MAIN_MIGRATE_CREATEFILE"
	envDbAuthMainMigrateInsertFile = "DB_AUTH_MAIN_MIGRATE_INSERTFILE"

	envDbBladeConnectHost     = "DB_BLADE_CONNECT_HOST"
	envDbBladeConnectPort     = "DB_BLADE_CONNECT_PORT"
	envDbBladeConnectSslMode  = "DB_BLADE_CONNECT_SSLMODE"
	envDbBladeConnectDbName   = "DB_BLADE_CONNECT_DBNAME"
	envDbBladeConnectUser     = "DB_BLADE_CONNECT_USER"
	envDbBladeConnectPassword = "DB_BLADE_CONNECT_PASSWORD"

	envDbBladeMigrateDropFile   = "DB_BLADE_MIGRATE_DROPFILE"
	envDbBladeMigrateCreateFile = "DB_BLADE_MIGRATE_CREATEFILE"

	envSmtpUser     = "SMTP_USER"
	envSmtpPassword = "SMTP_PASSWORD"
	envSmtpHost     = "SMTP_HOST"
	envSmtpPort     = "SMTP_PORT"

	envMailMessageFromName = "MAIL_MESSAGE_FROM_NAME"
	envMailMessageFromAddr = "MAIL_MESSAGE_FROM_ADDR"

	envCryptoSalt = "CRYPTO_SALT"

	envJwtIssuer                    = "JWT_ISSUER"
	envJwtAccessSecretKey           = "JWT_ACCESS_SECRET_KEY"
	envJwtAccessSignatureAlgorithm  = "JWT_ACCESS_SIGNATURE_ALGORITHM"
	envJwtAccessTokenLifeTime       = "JWT_ACCESS_TOKEN_LIFETIME"
	envJwtRefreshSecretKey          = "JWT_REFRESH_SECRET_KEY"
	envJwtRefreshSignatureAlgorithm = "JWT_REFRESH_SIGNATURE_ALGORITHM"
	envJwtRefreshTokenLifeTime      = "JWT_REFRESH_TOKEN_LIFETIME"
)
