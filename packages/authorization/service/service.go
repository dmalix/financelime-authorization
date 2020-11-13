/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/packages/authorization"
)

type Config struct {
	DomainAPP              string
	DomainAPI              string
	AuthInviteCodeRequired bool
	CryptoSalt             string
}

type Service struct {
	config          Config
	languageContent models.LanguageContent
	messageQueue    chan models.EmailMessage
	message         authorization.Message
	repository      authorization.Repository
	jwt             authorization.Jwt
}

func NewService(
	config Config,
	languageContent models.LanguageContent,
	messageQueue chan models.EmailMessage,
	message authorization.Message,
	repository authorization.Repository,
	jwt authorization.Jwt) *Service {
	return &Service{
		config:          config,
		languageContent: languageContent,
		messageQueue:    messageQueue,
		message:         message,
		repository:      repository,
		jwt:             jwt,
	}
}
