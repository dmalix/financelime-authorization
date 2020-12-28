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

type incomingProps struct {
	Email      string
	InviteCode string
	Language   string
}

func TestSignUp(t *testing.T) {

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              models.LanguageContent
		emailMessageQueue            = make(chan models.EmailMessage, 1)
		userMessage                  = new(email.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(repository.MockDescription)
		err                          error
		props                        incomingProps
		remoteAdrr                   = "127.0.0.1"
	)

	repository.MockData.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Confirm.Subject = append(languageContent.Data.User.Signup.Email.Confirm.Subject, "")
	languageContent.Data.User.Signup.Email.Confirm.Body = append(languageContent.Data.User.Signup.Email.Confirm.Body, "%s%s")

	cryptographerManager := cryptographer.NewCryptographer("6368616e676520746869732070617373")

	jwtManager := jwt.NewToken(
		"",
		"",
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

	err = newService.SignUp(props.Email, props.InviteCode, props.Language, remoteAdrr)

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}
