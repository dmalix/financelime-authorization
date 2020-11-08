/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"errors"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/packages/authorization/repository"
	"github.com/dmalix/financelime-rest-api/utils/email"
	"testing"
)

func TestConfirmUserEmail_Success(t *testing.T) {

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              models.LanguageContent
		emailMessageQueue            = make(chan models.EmailMessage, 1)
		userMessage                  = new(email.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(repository.MockDescription)
		err                          error
		message                      string
	)

	repository.MockData.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	//noinspection GoBoolExpressions
	var newService = NewService(
		configDomainAPI,
		configAuthInviteCodeRequired,
		languageContent,
		emailMessageQueue,
		userMessage,
		userRepo)

	message, err = newService.ConfirmUserEmail("12345")

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}

	if message != "text" {
		t.Errorf("service returned wrong the message value: got %v want %v",
			message, "text")
	}
}

func TestConfirmUserEmail_Error(t *testing.T) {

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
	email.AddEmailMessageToQueue_MockData.Expected.Error = errors.New("REPO_ERROR")

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	//noinspection GoBoolExpressions
	var newService = NewService(
		configDomainAPI,
		configAuthInviteCodeRequired,
		languageContent,
		emailMessageQueue,
		userMessage,
		userRepo)

	_, err = newService.ConfirmUserEmail("12345")

	if err == nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}


