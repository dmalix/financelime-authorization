/* Copyright © 2021. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"context"
	"errors"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"github.com/dmalix/financelime-authorization/app/authorization/repository"
	"github.com/dmalix/financelime-authorization/config"
	"github.com/dmalix/jwt"
	"github.com/dmalix/middleware"
	"github.com/dmalix/secretdata"
	"github.com/dmalix/sendmail"
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
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessageManager          = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		props                        incomingProps
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Request.Subject = append(languageContent.Data.User.Signup.Email.Request.Subject, "")
	languageContent.Data.User.Signup.Email.Request.Body = append(languageContent.Data.User.Signup.Email.Request.Body, "%s%s")

	cryptManager := secretdata.NewSecretData("6368616e676520746869732070617373")
	jwtManager := new(jwt.MockDescription)
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
		cryptManager,
		jwtManager,
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
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		message                      string
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	cryptographerManager := &secretdata.Cipher{}
	jwtManager := new(jwt.MockDescription)
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
		cryptographerManager,
		jwtManager,
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
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = errors.New("REPO_ERROR")

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Signup.Email.Password.Subject = append(languageContent.Data.User.Signup.Email.Password.Subject, "subject")
	languageContent.Data.User.Signup.Email.Password.Body = append(languageContent.Data.User.Signup.Email.Password.Body, "%s%s")
	languageContent.Data.User.Signup.Page.Text = append(languageContent.Data.User.Signup.Page.Text, "text")

	cryptographerManager := &secretdata.Cipher{}
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
		cryptographerManager,
		jwtManager,
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
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		device                       model.Device
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	tokenData := &secretdata.Cipher{}
	token := new(jwt.MockDescription)
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
		tokenData,
		tokenData,
		token,
		token)

	_, err = newService.CreateAccessToken(ctx, logger, model.ServiceCreateAccessTokenParam{
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
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(secretdata.MockDescription)
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
		cryptographerManager,
		jwtManager,
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
		accessTokenData              = []byte{123, 34, 73, 68, 34, 58, 50, 44, 34, 69, 109, 97, 105, 108, 34, 58, 34, 116, 101, 115, 116, 46, 117, 115, 101, 114, 64, 102, 105, 110, 97, 110, 99, 101, 108, 105, 109, 101, 46, 99, 111, 109, 34, 44, 34, 80, 97, 115, 115, 119, 111, 114, 100, 34, 58, 34, 34, 44, 34, 76, 97, 110, 103, 117, 97, 103, 101, 34, 58, 34, 101, 110, 34, 125}
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := new(secretdata.MockDescription)
	token := new(jwt.MockDescription)
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
		cryptographerManager,
		token,
		token)

	err = newService.RevokeRefreshToken(ctx, logger, model.ServiceRevokeRefreshTokenParam{
		AccessTokenData: accessTokenData,
		PublicSessionID: "publicSessionID"})

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
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.ResetPassword.Email.Request.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.ResetPassword.Email.Request.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	cryptographerManager := secretdata.NewSecretData("")
	token := new(jwt.MockDescription)
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
		cryptographerManager,
		token,
		token)

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
		accessTokenData              = []byte{123, 34, 73, 68, 34, 58, 50, 44, 34, 69, 109, 97, 105, 108, 34, 58, 34, 116, 101, 115, 116, 46, 117, 115, 101, 114, 64, 102, 105, 110, 97, 110, 99, 101, 108, 105, 109, 101, 46, 99, 111, 109, 34, 44, 34, 80, 97, 115, 115, 119, 111, 114, 100, 34, 58, 34, 34, 44, 34, 76, 97, 110, 103, 117, 97, 103, 101, 34, 58, 34, 101, 110, 34, 125}
		configDomainAPI              = "domain.com"
		configAuthInviteCodeRequired = true
		languageContent              config.LanguageContent
		emailMessageQueue            = make(chan sendmail.MessageBox, 1)
		emailMessage                 = new(sendmail.MockDescription)
		authRepo                     = new(repository.Mock)
		err                          error
		contextGetter                = new(middleware.MockDescription)
	)

	authRepo.Expected.Error = nil
	sendmail.MockData.Expected.Error = nil

	languageContent.Language = make(map[string]int)
	languageContent.Language["abc"] = 0
	languageContent.Data.User.Login.Email.Subject = append(languageContent.Data.User.Login.Email.Subject, "subject")
	languageContent.Data.User.Login.Email.Body = append(languageContent.Data.User.Login.Email.Body, "%s%s")

	secretData := new(secretdata.MockDescription)
	token := new(jwt.MockDescription)

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
		secretData,
		secretData,
		token,
		token)

	_, err = newService.GetListActiveSessions(ctx, logger, accessTokenData)
	if err != nil {
		t.Errorf("service returned wrong err value: got %v want %v",
			err, nil)
	}
}
