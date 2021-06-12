package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dmalix/authorization-service/app/authorization"
	"github.com/dmalix/authorization-service/app/authorization/model"
	"github.com/dmalix/authorization-service/config"
	"github.com/dmalix/jwt"
	"github.com/dmalix/middleware"
	"github.com/dmalix/requestid"
	"github.com/dmalix/secretdata"
	"github.com/dmalix/sendmail"
	"github.com/dmalix/utils/generate"
	"go.uber.org/zap"
	"net/mail"
	"time"
)

type service struct {
	config          model.ConfigService
	contextGetter   middleware.ContextGetter
	languageContent config.LanguageContent
	sendmailQueue   chan sendmail.MessageBox
	sendmailManager sendmail.Manager
	repository      authorization.Repository
	dataAccess      secretdata.SecretData
	dataRefresh     secretdata.SecretData
	jwtAccess       jwt.Jwt
	jwtRefresh      jwt.Jwt
}

func NewService(
	config model.ConfigService,
	contextGetter middleware.ContextGetter,
	languageContent config.LanguageContent,
	sendmailQueue chan sendmail.MessageBox,
	sendmailManager sendmail.Manager,
	repository authorization.Repository,
	dataAccess secretdata.SecretData,
	dataRefresh secretdata.SecretData,
	jwtAccess jwt.Jwt,
	jwtRefresh jwt.Jwt) *service {
	return &service{
		config:          config,
		contextGetter:   contextGetter,
		languageContent: languageContent,
		sendmailQueue:   sendmailQueue,
		sendmailManager: sendmailManager,
		repository:      repository,
		dataAccess:      dataAccess,
		dataRefresh:     dataRefresh,
		jwtAccess:       jwtAccess,
		jwtRefresh:      jwtRefresh,
	}
}

func (s *service) SignUpStep1(ctx context.Context, logger *zap.Logger, param model.ServiceSignUpParam) error {

	remoteAddr, remoteAddrKey, err := s.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get RemoteAddr", zap.Error(err))
		return err
	}

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	confirmationKey := generate.StringRand(16, 16, true)

	err = s.repository.SignUpStep1(ctx, logger, model.RepoSignUpParam{
		Email:              param.Email,
		Language:           param.Language,
		InviteCode:         param.InviteCode,
		ConfirmationKey:    confirmationKey,
		InviteCodeRequired: s.config.AuthInviteCodeRequired})
	if err != nil {
		logger.Error("failed to create new user", zap.Error(err), zap.String(requestIDKey, requestID))
		switch err {
		case authorization.ErrorBadParamEmail, authorization.ErrorBadParamInvite, authorization.ErrorBadParamLang,
			authorization.ErrorBadParamConfirmationKey:
			return authorization.ErrorBadParams
		case authorization.ErrorUserAlreadyExist:
			return err
		case authorization.ErrorInviteNotFound:
			return err
		case authorization.ErrorInviteHasEnded:
			return err
		default:
			return err
		}
	}

	newRequestID, err := requestid.Create(false)
	if err != nil {
		logger.DPanic("failed to generate requestID", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = s.sendmailManager.AddMessageToQueue(
		s.sendmailQueue,
		sendmail.Request{
			RemoteAddr:    remoteAddr,
			RemoteAddrKey: remoteAddrKey,
			RequestID:     requestID,
			RequestIDKey:  requestIDKey},
		sendmail.Email{
			To:      mail.Address{Address: param.Email},
			Subject: s.languageContent.Data.User.Signup.Email.Request.Subject[s.languageContent.Language[param.Language]],
			Body: fmt.Sprintf(
				s.languageContent.Data.User.Signup.Email.Request.Body[s.languageContent.Language[param.Language]],
				s.config.DomainAPI, confirmationKey, newRequestID),
			MessageID: fmt.Sprintf(
				"<%s@%s>",
				confirmationKey,
				fmt.Sprintf("%s.%s", "sign-up", s.config.DomainAPI))})
	if err != nil {
		logger.DPanic("failed to add an email message to the queue", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (s *service) SignUpStep2(ctx context.Context, logger *zap.Logger, confirmationKey string) (string, error) {

	remoteAddr, remoteAddrKey, err := s.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get RemoteAddr", zap.Error(err))
		return "", err
	}

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return "", err
	}

	user, err := s.repository.SignUpStep2(ctx, logger, confirmationKey)
	if err != nil {
		logger.Error("failed to confirm user email", zap.String(requestIDKey, requestID), zap.Error(err))
		switch err {
		case authorization.ErrorBadParamConfirmationKey:
			return "", err
		case authorization.ErrorConfirmationKeyNotFound, authorization.ErrorConfirmationKeyAlreadyConfirmed:
			return "", authorization.ErrorBadConfirmationKey
		default:
			return "", err
		}
	}

	err = s.sendmailManager.AddMessageToQueue(
		s.sendmailQueue,
		sendmail.Request{
			RemoteAddr:    remoteAddr,
			RemoteAddrKey: remoteAddrKey,
			RequestID:     requestID,
			RequestIDKey:  requestIDKey,
		},
		sendmail.Email{
			To:      mail.Address{Address: user.Email},
			Subject: s.languageContent.Data.User.Signup.Email.Password.Subject[s.languageContent.Language[user.Language]],
			Body: fmt.Sprintf(
				s.languageContent.Data.User.Signup.Email.Password.Body[s.languageContent.Language[user.Language]],
				user.Password),
			MessageID: fmt.Sprintf(
				"<%s@%s>",
				user.Password,
				fmt.Sprintf("%s.%s", "confirm-user-email", s.config.DomainAPI))})
	if err != nil {
		logger.DPanic("failed to send message to the user", zap.String(requestIDKey, requestID), zap.Error(err))
		return "", err
	}

	confirmationMessage := s.languageContent.Data.User.Signup.Page.Text[s.languageContent.Language[user.Language]]

	return confirmationMessage, nil
}

func (s *service) CreateAccessToken(ctx context.Context, logger *zap.Logger,
	param model.ServiceCreateAccessTokenParam) (model.ServiceAccessTokenReturn, error) {

	remoteAddr, remoteAddrKey, err := s.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return model.ServiceAccessTokenReturn{}, err
	}
	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.ServiceAccessTokenReturn{}, err
	}

	user, err := s.repository.GetUserByAuth(ctx, logger, model.RepoGetUserByAuthParam{
		Email:    param.Email,
		Password: param.Password,
	})
	if err != nil {
		logger.Error("failed to create access token", zap.Error(err), zap.String(requestIDKey, requestID))
		switch err {
		case authorization.ErrorBadParamEmail, authorization.ErrorBadParamPassword, authorization.ErrorBadParamLang:
			return model.ServiceAccessTokenReturn{}, authorization.ErrorBadParams
		case authorization.ErrorUserNotFound:
			return model.ServiceAccessTokenReturn{}, err
		default:
			return model.ServiceAccessTokenReturn{}, err
		}
	}

	publicSessionID, err := generate.PublicID(user.ID)
	if err != nil {
		logger.DPanic("failed to generate the publicSessionID", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}
	userData, err := json.Marshal(user)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}
	encryptedAccessTokenData, err := s.dataAccess.Encrypt(userData)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	encryptedRefreshTokenData, err := s.dataRefresh.Encrypt(userData)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	accessToken, err := s.jwtAccess.Create(jwt.Claims{
		JwtID: publicSessionID,
		Data:  encryptedAccessTokenData,
	})
	if err != nil {
		logger.DPanic("failed to generate an access token", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	refreshToken, err := s.jwtRefresh.Create(jwt.Claims{
		JwtID: publicSessionID,
		Data:  encryptedRefreshTokenData,
	})
	if err != nil {
		logger.DPanic("failed to generate an refresh token", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	err = s.repository.CreateSession(ctx, logger, model.RepoCreateSessionParam{
		UserID:          user.ID,
		PublicSessionID: publicSessionID,
		RefreshToken:    refreshToken,
		ClientID:        param.ClientID,
		UserAgent:       param.UserAgent,
		Device:          param.Device})
	if err != nil {
		logger.DPanic("failed to create session", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	err = s.sendmailManager.AddMessageToQueue(
		s.sendmailQueue,
		sendmail.Request{
			RemoteAddr:    remoteAddr,
			RemoteAddrKey: remoteAddrKey,
			RequestID:     requestID,
			RequestIDKey:  requestIDKey},
		sendmail.Email{
			To:      mail.Address{Address: param.Email},
			Subject: s.languageContent.Data.User.Login.Email.Subject[s.languageContent.Language[user.Language]],
			Body: fmt.Sprintf(
				s.languageContent.Data.User.Login.Email.Body[s.languageContent.Language[user.Language]],
				time.Now().UTC().String(),
				param.Device.Platform,
				remoteAddr,
				s.config.DomainAPP),
			MessageID: fmt.Sprintf(
				"<%s@%s>",
				remoteAddr,
				fmt.Sprintf("%s.%s", "get-access-token", s.config.DomainAPI))})
	if err != nil {
		logger.DPanic("Failed to add this email message to the queue", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	return model.ServiceAccessTokenReturn{
		PublicSessionID: publicSessionID,
		AccessJWT:       accessToken,
		RefreshJWT:      refreshToken}, nil
}

func (s *service) RefreshAccessToken(ctx context.Context, logger *zap.Logger,
	refreshToken string) (model.ServiceAccessTokenReturn, error) {

	var JwtAccessData []byte
	var JwtRefreshData []byte

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.ServiceAccessTokenReturn{}, err
	}

	jwtSource, parseCodeError, err := s.jwtRefresh.Parse(refreshToken)
	if err != nil {
		logger.Error("failed to verify the refresh token", zap.Error(err),
			zap.String("refreshToken", refreshToken),
			zap.String("parseCodeError", parseCodeError),
			zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, authorization.ErrorBadRefreshToken
	}

	user, err := s.repository.GetUserByRefreshToken(ctx, logger, refreshToken)
	logger.Error("failed to refresh the access token", zap.Error(err), zap.String(requestIDKey, requestID))
	if err != nil {
		switch err {
		case authorization.ErrorUserNotFound:
			return model.ServiceAccessTokenReturn{}, err
		default:
			return model.ServiceAccessTokenReturn{}, err
		}
	}

	publicSessionID := jwtSource.Claims.JwtID

	sourceUserData, err := json.Marshal(user)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	if s.config.JwtAccessEncryptData {
		JwtAccessData, err = s.dataAccess.Encrypt(sourceUserData)
		if err != nil {
			logger.DPanic("failed to encrypt the access data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.ServiceAccessTokenReturn{}, err
		}
	} else {
		JwtAccessData = sourceUserData
	}

	if s.config.JwtRefreshEncryptData {
		JwtRefreshData, err = s.dataRefresh.Encrypt(sourceUserData)
		if err != nil {
			logger.DPanic("failed to encrypt the refresh data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.ServiceAccessTokenReturn{}, err
		}
	} else {
		JwtRefreshData = sourceUserData
	}

	jwtAccess, err := s.jwtAccess.Create(jwt.Claims{
		JwtID: publicSessionID,
		Data:  JwtAccessData,
	})
	if err != nil {
		logger.DPanic("failed to create an access token (JWT)", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	jwtRefresh, err := s.jwtRefresh.Create(jwt.Claims{
		JwtID: publicSessionID,
		Data:  JwtRefreshData,
	})
	if err != nil {
		logger.DPanic("failed to create an refresh token (JWT)", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	err = s.repository.UpdateSession(ctx, logger, model.RepoUpdateSessionParam{
		PublicSessionID: publicSessionID,
		RefreshToken:    refreshToken,
	})
	if err != nil {
		logger.DPanic("failed to update the session", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	return model.ServiceAccessTokenReturn{
		PublicSessionID: publicSessionID,
		AccessJWT:       jwtAccess,
		RefreshJWT:      jwtRefresh}, nil
}

func (s *service) RevokeRefreshToken(ctx context.Context, logger *zap.Logger, param model.ServiceRevokeRefreshTokenParam) error {

	var user model.User

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	err = json.Unmarshal(param.AccessTokenData, &user)
	if err != nil {
		logger.DPanic("failed to unmarshal the decryptedJWTData value", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = s.repository.DeleteSession(ctx, logger, model.RepoDeleteSessionParam{
		UserID:          user.ID,
		PublicSessionID: param.PublicSessionID})
	if err != nil {
		logger.DPanic("failed to delete the session", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (s *service) GetListActiveSessions(ctx context.Context, logger *zap.Logger, accessTokenData []byte) ([]model.Session, error) {

	var user model.User

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return nil, err
	}

	err = json.Unmarshal(accessTokenData, &user)
	if err != nil {
		logger.DPanic("failed to unmarshal the accessTokenData value to struct", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return nil, err
	}

	sessions, err := s.repository.GetListActiveSessions(ctx, logger, user.ID)
	if err != nil {
		logger.DPanic("failed to get the active sessions list", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return nil, err
	}
	return sessions, nil
}

func (s *service) ResetUserPasswordStep1(ctx context.Context, logger *zap.Logger, email string) error {

	remoteAddr, remoteAddrKey, err := s.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return err
	}

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	confirmationKey := generate.StringRand(16, 16, true)

	user, err := s.repository.ResetUserPasswordStep1(ctx, logger, model.RepoResetUserPasswordParam{
		Email:           email,
		ConfirmationKey: confirmationKey})
	if err != nil {
		logger.Error("failed to request a reset of the user's password", zap.Error(err), zap.String(requestIDKey, requestID))
		switch err {
		case authorization.ErrorBadParamEmail, authorization.ErrorBadConfirmationKey:
			return err
		case authorization.ErrorUserNotFound:
			return err
		default:
			return err
		}
	}

	newRequestID, err := requestid.Create(false)
	if err != nil {
		logger.DPanic("failed to generate requestID", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = s.sendmailManager.AddMessageToQueue(
		s.sendmailQueue,
		sendmail.Request{
			RemoteAddr:    remoteAddr,
			RemoteAddrKey: remoteAddrKey,
			RequestID:     requestID,
			RequestIDKey:  requestIDKey},
		sendmail.Email{
			To:      mail.Address{Address: email},
			Subject: s.languageContent.Data.User.ResetPassword.Email.Request.Subject[s.languageContent.Language[user.Language]],
			Body: fmt.Sprintf(
				s.languageContent.Data.User.ResetPassword.Email.Request.Body[s.languageContent.Language[user.Language]],
				remoteAddr, s.config.DomainAPI, confirmationKey, newRequestID),
			MessageID: fmt.Sprintf(
				"<%s@%s>",
				confirmationKey,
				fmt.Sprintf("%s.%s", "reset-password", s.config.DomainAPI))})
	if err != nil {
		logger.DPanic("failed to add an email message to the queue", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (s *service) ResetUserPasswordStep2(ctx context.Context, logger *zap.Logger, confirmationKey string) (string, error) {

	remoteAddr, remoteAddrKey, err := s.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get RemoteAddr", zap.Error(err))
		return "", err
	}

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return "", err
	}

	user, err := s.repository.ResetUserPasswordStep2(ctx, logger, confirmationKey)
	if err != nil {
		logger.Error("failed to confirm user password reset", zap.String(requestIDKey, requestID), zap.Error(err))
		switch err {
		case authorization.ErrorBadParamConfirmationKey:
			return "", err
		case authorization.ErrorConfirmationKeyNotFound, authorization.ErrorConfirmationKeyAlreadyConfirmed,
			authorization.ErrorUserNotFound:
			return "", authorization.ErrorBadConfirmationKey
		default:
			return "", err
		}
	}

	err = s.sendmailManager.AddMessageToQueue(
		s.sendmailQueue,
		sendmail.Request{
			RemoteAddr:    remoteAddr,
			RemoteAddrKey: remoteAddrKey,
			RequestID:     requestID,
			RequestIDKey:  requestIDKey,
		},
		sendmail.Email{
			To:      mail.Address{Address: user.Email},
			Subject: s.languageContent.Data.User.ResetPassword.Email.Password.Subject[s.languageContent.Language[user.Language]],
			Body: fmt.Sprintf(
				s.languageContent.Data.User.ResetPassword.Email.Password.Body[s.languageContent.Language[user.Language]],
				user.Password),
			MessageID: fmt.Sprintf(
				"<%s@%s>",
				user.Password,
				fmt.Sprintf("%s.%s", "confirm-user-password-reset", s.config.DomainAPI))})
	if err != nil {
		logger.DPanic("failed to send message to the user", zap.String(requestIDKey, requestID), zap.Error(err))
		return "", err
	}

	confirmationMessage := s.languageContent.Data.User.ResetPassword.Page.Text[s.languageContent.Language[user.Language]]

	return confirmationMessage, nil
}
