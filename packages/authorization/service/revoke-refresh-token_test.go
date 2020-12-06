/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"github.com/dmalix/financelime-authorization/models"
	"github.com/dmalix/financelime-authorization/packages/authorization/repository"
	"github.com/dmalix/financelime-authorization/utils/cryptographer"
	"github.com/dmalix/financelime-authorization/utils/email"
	"github.com/dmalix/financelime-authorization/utils/jwt"
	"testing"
)

func TestRevokeRefreshToken(t *testing.T) {

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              models.LanguageContent
		emailMessageQueue            = make(chan models.EmailMessage, 1)
		userMessage                  = new(email.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(repository.MockDescription)
		err                          error
	)

	repository.MockData.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer.MockDescription)

	jwtManager := jwt.NewToken(
		"12345",
		jwt.PropsSigningAlgorithmHS256,
		"",
		"",
		0,
		0)

	serviceConfig := ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		languageContent,
		emailMessageQueue,
		userMessage,
		userRepo,
		cryptographerManager,
		jwtManager)

	err = newService.RevokeRefreshToken([]byte("encryptedUserData"), "publicSessionID")

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}
