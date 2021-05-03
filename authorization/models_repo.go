package authorization

import "database/sql"

type ConfigPostgreDB struct {
	Host     string
	Port     string
	SSLMode  string
	DBName   string
	User     string
	Password string
}

type ConfigRepository struct {
	CryptoSalt              string
	JwtRefreshTokenLifetime int
}

type repository struct {
	config     ConfigRepository
	dbAuthMain *sql.DB
	dbAuthRead *sql.DB
	dbBlade    *sql.DB
}

type repoSaveSessionParam struct {
	userID          int64
	publicSessionID string
	refreshToken    string
	clientID        string
	remoteAddr      string
	userAgent       string
	device          device
}

type repoCreateUserParam struct {
	email              string
	language           string
	inviteCode         string
	remoteAddr         string
	confirmationKey    string
	inviteCodeRequired bool
}

type repoGetUserByAuthParam struct {
	email    string
	password string
}

type repoDeleteSessionParam struct {
	userID          int64
	publicSessionID string
}

type repoUpdateSessionParam struct {
	publicSessionID string
	refreshToken    string
	remoteAddr      string
}

type repoRequestUserPasswordResetParam struct {
	email           string
	remoteAddr      string
	confirmationKey string
}
