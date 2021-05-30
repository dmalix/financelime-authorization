/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"context"
	"errors"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"github.com/dmalix/financelime-authorization/app/authorization/repository"
	"github.com/dmalix/financelime-authorization/config"
	"github.com/dmalix/financelime-authorization/packages/cryptographer"
	"github.com/dmalix/financelime-authorization/packages/email"
	"github.com/dmalix/financelime-authorization/packages/jwt"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"go.uber.org/zap"
	"testing"
)

const remoteAddr = "127.0.0.1"
const requestID = "W7000-T6755-T7700-P4010-W6778"

func TestServiceSignUp(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, middleware.ContextKeyRemoteAddr, remoteAddr)
	ctx = context.WithValue(ctx, middleware.ContextKeyRequestID, requestID)
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessageManager          = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		props                        incomingProps
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Request.Subject = append(languageContent.Data.User.Signup.Email.Request.Subject, "")
	languageContent.Data.User.Signup.Email.Request.Body = append(languageContent.Data.User.Signup.Email.Request.Body, "%s%s")

	cryptManager := cryptographer.NewCryptographer("6368616e676520746869732070617373")
	jwtManager := &jwt.Token{}
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	contextGetter := middleware.NewContextGetter()

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessageManager,
		authRepo,
		cryptManager,
		jwtManager)

	err = newService.SignUpStep1(ctx, logger, model.ServiceSignUpParam{
		Email:      props.Email,
		Language:   props.Language,
		InviteCode: props.InviteCode,
	})

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}

func TestServiceConfirmUserEmail_Success(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		message                      string
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	cryptographerManager := &cryptographer.Cipher{}
	jwtManager := &jwt.Token{}
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	message, err = newService.SignUpStep2(ctx, logger, "12345")

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

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = errors.New("REPO_ERROR")

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	cryptographerManager := &cryptographer.Cipher{}
	jwtManager := &jwt.Token{}
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	_, err = newService.SignUpStep2(ctx, logger, "12345")

	if err == nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRequestAccessToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		device                       model.Device
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := &cryptographer.Cipher{}
	jwtManager := jwt.NewToken("12345", jwt.ParamSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	_, err =
		newService.CreateAccessToken(ctx, logger, model.ServiceCreateAccessTokenParam{
			Email:     "email",
			Password:  "password",
			ClientID:  "PWA",
			UserAgent: "userAgent",
			Device:    device,
		})

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRefreshAccessToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer.MockDescription)
	jwtManager := new(jwt.MockDescription)
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	_, err = newService.RefreshAccessToken(ctx, logger, "refreshToken")

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRevokeRefreshToken(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer.MockDescription)
	jwtManager := jwt.NewToken("12345", jwt.ParamSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	err = newService.RevokeRefreshToken(ctx, logger, model.ServiceRevokeRefreshTokenParam{
		EncryptedUserData: []byte("encryptedUserData"),
		PublicSessionID:   "publicSessionID"})

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceRequestUserPasswordReset(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.ResetPassword.Email.Request.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.ResetPassword.Email.Request.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := cryptographer.NewCryptographer("")
	jwtManager := jwt.NewToken("12345", jwt.ParamSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	err = newService.ResetUserPasswordStep1(ctx, logger, "email")

	if err != nil {
		t.Errorf("service returned wrong the err value: got %v want %v",
			err, nil)
	}
}

func TestServiceGetListActiveSessions(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewProduction()

	var (
		encryptedUserData            []byte
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan email.MessageBox, 1)
		emailMessage                 = new(email.AddEmailMessageToQueue_MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	email.AddEmailMessageToQueue_MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(cryptographer.MockDescription)
	jwtManager := jwt.NewToken("12345", jwt.ParamSigningAlgorithmHS256, "", "", 0, 0)
	//goland:noinspection GoBoolExpressions
	serviceConfig := model.ConfigService{
		DomainAPI:              configDomainAPI,
		AuthInviteCodeRequired: configAuthInviteCodeRequired,
	}

	//noinspection GoBoolExpressions
	var newService = NewService(
		serviceConfig,
		contextGetter,
		languageContent,
		emailMessageQueue,
		emailMessage,
		authRepo,
		cryptographerManager,
		jwtManager)

	_, err = newService.GetListActiveSessions(ctx, logger, encryptedUserData)

	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}
