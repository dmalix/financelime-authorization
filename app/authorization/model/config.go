package model

type ConfigService struct {
	DomainAPP              string
	DomainAPI              string
	AuthInviteCodeRequired bool
	SecretKey              string
	CryptoSalt             string
}

type ConfigRepository struct {
	CryptoSalt              string
	JwtRefreshTokenLifetime int
}

type ConfigPostgresDB struct {
	Host     string
	Port     int
	SSLMode  string
	DBName   string
	User     string
	Password string
}
