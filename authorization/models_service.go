package authorization

import (
	cryptographer3 "github.com/dmalix/financelime-authorization/packages/cryptographer"
	email2 "github.com/dmalix/financelime-authorization/packages/email"
	jwt3 "github.com/dmalix/financelime-authorization/packages/jwt"
)

type service struct {
	config          ConfigService
	languageContent LanguageContent
	messageQueue    chan email2.EmailMessage
	message         email2.Message
	repository      Repository
	cryptographer   cryptographer3.Cryptographer
	jwt             jwt3.Jwt
}

type ConfigService struct {
	DomainAPP              string
	DomainAPI              string
	AuthInviteCodeRequired bool
	SecretKey              string
	CryptoSalt             string
}

type serviceSignUpParam struct {
	email      string
	language   string
	inviteCode string
	remoteAddr string
}

type serviceCreateAccessTokenParam struct {
	email      string
	password   string
	clientID   string
	remoteAddr string
	userAgent  string
	device     device
}

type serviceRefreshAccessTokenParam struct {
	refreshToken string
	remoteAddr   string
}

type serviceAccessTokenReturn struct {
	publicSessionID string
	accessJWT       string
	refreshJWT      string
}

type serviceRevokeRefreshTokenParam struct {
	encryptedUserData []byte
	publicSessionID   string
}

type serviceRequestUserPasswordResetParam struct {
	email      string
	remoteAddr string
}
