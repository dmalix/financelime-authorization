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
	LanguageContent struct {
		File string
	}
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
	Crypto struct {
		Salt string
	}
	Jwt struct {
		SecretKey            string
		SigningAlgorithm     string
		Issuer               string
		Subject              string
		AccessTokenLifetime  int
		RefreshTokenLifetime int
	}
}
