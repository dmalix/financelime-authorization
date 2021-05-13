/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"context"
	"errors"
	cryptographer2 "github.com/dmalix/financelime-authorization/packages/cryptographer"
	email2 "github.com/dmalix/financelime-authorization/packages/email"
	jwt2 "github.com/dmalix/financelime-authorization/packages/jwt"
	"testing"
)

func TestServiceSignUp(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
		props                        incomingProps
		remoteAddr                   = "127.0.0.1"
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Confirm.Subject = append(languageContent.Data.User.Signup.Email.Confirm.Subject, "")
	languageContent.Data.User.Signup.Email.Confirm.Body = append(languageContent.Data.User.Signup.Email.Confirm.Body, "%s%s")

	cryptographerManager := cryptographer2.NewCryptographer("6368616e676520746869732070617373")
	jwtManager := &jwt2.Token{}
	//goland:noinspection GoBoolExpressions
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

	err = newService.signUp(ctx, serviceSignUpParam{
		email:      props.Email,
		language:   props.Language,
		inviteCode: props.InviteCode,
		remoteAddr: remoteAddr,
	})

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}

func TestServiceConfirmUserEmail_Success(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
		message                      string
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	cryptographerManager := &cryptographer2.Cipher{}
	jwtManager := &jwt2.Token{}
	//goland:noinspection GoBoolExpressions
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

	message, err = newService.confirmUserEmail(ctx, "12345")

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}

	if message != "text" {
		t.Errorf("service returned wrong the message value: got %v want %v",
			message, "text")
	}
}

func TestServiceConfirmUserEmail_Error(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = errors.New("REPO_ERROR")

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	cryptographerManager := &cryptographer2.Cipher{}
	jwtManager := &jwt2.Token{}
	//goland:noinspection GoBoolExpressions
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

	_, err = newService.confirmUserEmail(ctx, "12345")

	if err == nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRequestAccessToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
		device                       device
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := &cryptographer2.Cipher{}
	jwtManager := jwt2.NewToken("12345", jwt2.PropsSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
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

	_, err =
		newService.createAccessToken(ctx, serviceCreateAccessTokenParam{
			email:      "email",
			password:   "password",
			clientID:   "PWA",
			remoteAddr: "127.0.0.1",
			userAgent:  "userAgent",
			device:     device,
		})

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRefreshAccessToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer2.MockDescription)
	jwtManager := new(jwt2.MockDescription)
	//goland:noinspection GoBoolExpressions
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

	_, err = newService.refreshAccessToken(ctx, serviceRefreshAccessTokenParam{
		refreshToken: "refreshToken",
		remoteAddr:   "127.0.0.1",
	})

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRevokeRefreshToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer2.MockDescription)
	jwtManager := jwt2.NewToken("12345", jwt2.PropsSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
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

	err = newService.revokeRefreshToken(ctx, serviceRevokeRefreshTokenParam{
		[]byte("encryptedUserData"),
		"publicSessionID"})

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRequestUserPasswordReset(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.ResetPassword.Email.Request.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.ResetPassword.Email.Request.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := cryptographer2.NewCryptographer("")
	jwtManager := jwt2.NewToken("12345", jwt2.PropsSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
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

	err = newService.requestUserPasswordReset(ctx, serviceRequestUserPasswordResetParam{"email", "127.0.0.1"})

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceGetListActiveSessions(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var encryptedUserData []byte

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              LanguageContent
		emailMessageQueue            = make(chan email2.EmailMessage, 1)
		userMessage                  = new(email2.AddEmailMessageToQueue_MockDescription)
		userRepo                     = new(RepoMockDescription)
		err                          error
	)

	ServiceMockData.Expected.Error = nil
	email2.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer2.MockDescription)
	jwtManager := jwt2.NewToken("12345", jwt2.PropsSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
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

	_, err = newService.getListActiveSessions(ctx, encryptedUserData)

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}
