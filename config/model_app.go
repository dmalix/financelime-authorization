package config

import (
	"net/mail"
)

type Version struct {
	DevelopmentMode bool
	Number          string
	BuildTime       string
	Commit          string
	Compiler        string
}

type App struct {
	Domain struct {
		App string
		Api string
	}
	Http struct {
		Port int
	}
	Auth struct {
		InviteCodeRequired bool
	}
	Db struct {
		AuthMain DB
		AuthRead DB
		Blade    DB
	}
	Smtp struct {
		User     string
		Password string
		Host     string
		Port     int
	}
	MailMessage struct {
		From mail.Address
	}
	LanguageContent struct {
		File string
	}
	Crypto struct {
		Salt string
	}
	Jwt struct {
		Issuer                    string
		AccessAudience            string
		AccessSecretKey           string
		AccessSignatureAlgorithm  string
		AccessEncryptData         bool
		AccessTokenLifetime       int
		RefreshAudience           string
		RefreshSecretKey          string
		RefreshSignatureAlgorithm string
		RefreshEncryptData        bool
		RefreshTokenLifetime      int
	}
}
