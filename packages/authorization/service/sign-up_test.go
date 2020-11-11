/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/packages/authorization/repository"
	"github.com/dmalix/financelime-rest-api/utils/email"
	"github.com/dmalix/financelime-rest-api/utils/jwt"
	"testing"
)

type incomingProps struct {
	Email      string
	InviteCode string
	Language   string
}

func TestSignUp_Success(t *testing.T) {

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

	jwtManager := jwt.NewToken(
		"",
		"",
		"",
		"",
		0,
		0)

	serviceConfig := Config{
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
		jwtManager)

	err = newService.SignUp(props.Email, props.InviteCode, props.Language, remoteAdrr)

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}

func TestSignUp_RepoError(t *testing.T) {

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              models.LanguageContent
		emailMessageQueue            = make(chan models.EmailMessage, 1)
		userMessage                  = new(email.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(repository.MockDescription)
		err                          error
		errDescription               = "4PtDRMCQ:a system error was returned[RepoError]"
		props                        incomingProps
		remoteAdrr                   = "127.0.0.1"
	)

	repository.MockData.Expected.Error = errors.New("RepoError")

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Confirm.Subject = append(languageContent.Data.User.Signup.Email.Confirm.Subject, "")
	languageContent.Data.User.Signup.Email.Confirm.Body = append(languageContent.Data.User.Signup.Email.Confirm.Body, "%s%s")

	jwtManager := jwt.NewToken(
		"",
		"",
		"",
		"",
		0,
		0)

	serviceConfig := Config{
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
		jwtManager)

	err = newService.SignUp(props.Email, props.InviteCode, props.Language, remoteAdrr)

	if err != nil {
		if err.Error() != errDescription {
			t.Errorf("service returned wrong err value: got %v want %v",
				err, errDescription)
		}
	} else {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, errDescription)
	}
}

func TestSignUp_EmailError(t *testing.T) {

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              models.LanguageContent
		emailMessageQueue            = make(chan models.EmailMessage, 1)
		userMessage                  = new(email.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(repository.MockDescription)
		err                          error
		errDescription               = "XfCCWkb2:Failed to send message to the user[EmailError]"
		props                        incomingProps
		remoteAdrr                   = "127.0.0.1"
	)

	repository.MockData.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = errors.New("EmailError")

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Confirm.Subject =
		append(languageContent.Data.User.Signup.Email.Confirm.Subject,
			"")
	languageContent.Data.User.Signup.Email.Confirm.Body =
		append(languageContent.Data.User.Signup.Email.Confirm.Body,
			"%s%s")

	jwtManager := jwt.NewToken(
		"",
		"",
		"",
		"",
		0,
		0)

	serviceConfig := Config{
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
		jwtManager)

	err = newService.SignUp(props.Email, props.InviteCode, props.Language, remoteAdrr)

	if err != nil {
		if err.Error() != errDescription {
			t.Errorf("service returned wrong err value: got %v want %v",
				err, errDescription)
		}
	} else {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, errDescription)
	}
}
