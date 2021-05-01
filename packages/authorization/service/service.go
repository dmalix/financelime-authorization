/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/packages/authorization"
)

type Service struct {
	domainAPI          string
	inviteCodeRequired bool
	languageContent    models.LanguageContent
	messageQueue       chan models.EmailMessage
	userMessage        authorization.UserMessage
	userRepo           authorization.UserRepo
}

func NewService(
	domainAPI string,
	inviteCodeRequired bool,
	languageContent models.LanguageContent,
	messageQueue chan models.EmailMessage,
	userMessage authorization.UserMessage,
	userRepo authorization.UserRepo) *Service {
	return &Service{
		domainAPI:          domainAPI,
		inviteCodeRequired: inviteCodeRequired,
		languageContent:    languageContent,
		messageQueue:       messageQueue,
		userMessage:        userMessage,
		userRepo:           userRepo,
	}
}
