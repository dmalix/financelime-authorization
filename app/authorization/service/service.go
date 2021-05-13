/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dmalix/financelime-authorization/app/authorization"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"github.com/dmalix/financelime-authorization/config"
	"github.com/dmalix/financelime-authorization/packages/cryptographer"
	"github.com/dmalix/financelime-authorization/packages/email"
	"github.com/dmalix/financelime-authorization/packages/jwt"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"github.com/dmalix/financelime-authorization/utils/random"
	"go.uber.org/zap"
	"net/mail"
	"strconv"
	"time"
)

type service struct {
	config          model.ConfigService
	contextGetter   middleware.ContextGetter
	languageContent config.LanguageContent
	messageQueue    chan email.EMessage
	message         email.Message
	repository      authorization.Repository
	cryptographer   cryptographer.Cryptographer
	jwt             jwt.JWT
}

func NewService(
	config model.ConfigService,
	contextGetter middleware.ContextGetter,
	languageContent config.LanguageContent,
	messageQueue chan email.EMessage,
	message email.Message,
	repository authorization.Repository,
	cryptographer cryptographer.Cryptographer,
	jwt jwt.JWT) *service {
	return &service{
		config:          config,
		contextGetter:   contextGetter,
		languageContent: languageContent,
		messageQueue:    messageQueue,
		message:         message,
		repository:      repository,
		cryptographer:   cryptographer,
		jwt:             jwt,
	}
}

func (s *service) generatePublicID(ctx context.Context, logger *zap.Logger, privateID int64) (string, error) {

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return "", err
	}

	hs := sha256.New()

	_, err = hs.Write([]byte(
		strconv.FormatInt(privateID, 10) +
			random.StringRand(16, 16, false) +
			time.Now().String() +
			s.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate a hash", zap.Error(err), zap.String(requestIDKey, requestID))
		return "", err
	}

	publicSessionID := hex.EncodeToString(hs.Sum(nil))

	return publicSessionID, nil
}

func (s *service) generateRequestID(_ context.Context) (string, error) {

	const requestID = "W7000-T6755-T7700-P4010-W6778"

	//TODO Add the generation requestID

	/*

			#### Algorithm for generating the 'request-id' header:

			```
		        +------------------------- | 4 random character sets
		        |     +------------------- |          pattern:
		        |     |     +------------- |            LNNNN
		        |     |     |     +------- |       (1 letter + 4 numbers)
		check   |     |     |     |
		 sum    1     2     3     4
		----- ----- ----- ----- -----
		LXXXX-LNNNN-LNNNN-LNNNN-LNNNN  <-- pattern
		W7000-T6755-T7700-P4010-W6778  <-- example of generated value
		=----                   =----
		^                       |
		|                       |
		+-----------------------+

		LXXXX | L - The letter that is always equal to a letter from set 4
		      | XXXX - This is the sum of all digits from sets 1,2,3 and 4
		      |        The number of characters is always 4 (missing ones get zeros)

	*/

	return requestID, nil
}

func (s *service) SignUp(ctx context.Context, logger *zap.Logger, param model.ServiceSignUpParam) error {

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	confirmationKey := random.StringRand(16, 16, true)

	err = s.repository.CreateUser(ctx, logger, model.RepoCreateUserParam{
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

	newRequestID, err := s.generateRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to generate requestID", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: param.Email},
		s.languageContent.Data.User.Signup.Email.Confirm.Subject[s.languageContent.Language[param.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Signup.Email.Confirm.Body[s.languageContent.Language[param.Language]],
			s.config.DomainAPI, confirmationKey, newRequestID),
		fmt.Sprintf(
			"<%s@%s>",
			confirmationKey,
			fmt.Sprintf("%s.%s", "sign-up", s.config.DomainAPI)))
	if err != nil {
		logger.DPanic("failed to add an email message to the queue", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (s *service) ConfirmUserEmail(ctx context.Context, logger *zap.Logger, confirmationKey string) (string, error) {
	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return "", err
	}

	user, err := s.repository.ConfirmUserEmail(ctx, logger, confirmationKey)
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

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: user.Email},
		s.languageContent.Data.User.Signup.Email.Password.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Signup.Email.Password.Body[s.languageContent.Language[user.Language]],
			user.Password),
		fmt.Sprintf(
			"<%s@%s>",
			user.Password,
			fmt.Sprintf("%s.%s", "confirm-user-email", s.config.DomainAPI)))
	if err != nil {
		logger.DPanic("failed to send message to the user", zap.String(requestIDKey, requestID), zap.Error(err))
		return "", err
	}

	confirmationMessage := s.languageContent.Data.User.Signup.Page.Text[s.languageContent.Language[user.Language]]

	return confirmationMessage, nil
}

func (s *service) CreateAccessToken(ctx context.Context, logger *zap.Logger,
	param model.ServiceCreateAccessTokenParam) (model.ServiceAccessTokenReturn, error) {

	remoteAddr, err := s.contextGetter.GetRemoteAddr(ctx)
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

	publicSessionID, err := s.generatePublicID(ctx, logger, user.ID)
	if err != nil {
		logger.DPanic("failed to generate the publicSessionID", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	sourceUserData, err := json.Marshal(user)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	encryptedUserData, err := s.cryptographer.Encrypt(sourceUserData)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	accessJWT, err := s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.ParamPurposeAccess)
	if err != nil {
		logger.DPanic("failed to generate an access token", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	refreshJWT, err := s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.ParamPurposeRefresh)
	if err != nil {
		logger.DPanic("failed to generate an refresh token", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	err = s.repository.SaveSession(ctx, logger, model.RepoSaveSessionParam{
		UserID:          user.ID,
		PublicSessionID: publicSessionID,
		RefreshToken:    refreshJWT,
		ClientID:        param.ClientID,
		UserAgent:       param.UserAgent,
		Device:          param.Device})
	if err != nil {
		logger.DPanic("failed to save session", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: param.Email},
		s.languageContent.Data.User.Login.Email.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Login.Email.Body[s.languageContent.Language[user.Language]],
			time.Now().UTC().String(),
			param.Device.Platform,
			remoteAddr,
			s.config.DomainAPP),
		fmt.Sprintf(
			"<%s@%s>",
			remoteAddr,
			fmt.Sprintf("%s.%s", "get-access-token", s.config.DomainAPI)))
	if err != nil {
		logger.DPanic("Failed to add this email message to the queue", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	return model.ServiceAccessTokenReturn{
		PublicSessionID: publicSessionID,
		AccessJWT:       accessJWT,
		RefreshJWT:      refreshJWT}, nil
}

func (s *service) RefreshAccessToken(ctx context.Context, logger *zap.Logger,
	refreshToken string) (model.ServiceAccessTokenReturn, error) {

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.ServiceAccessTokenReturn{}, err
	}

	jwtData, err := s.jwt.VerifyToken(refreshToken)
	if err != nil {
		logger.Error("failed to verify the refresh token", zap.Error(err),
			zap.String("refreshToken", refreshToken),
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

	publicSessionID := jwtData.Payload.PublicSessionID

	sourceUserData, err := json.Marshal(user)
	if err != nil {
		logger.DPanic("failed to marshal the user struct", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	encryptedUserData, err := s.cryptographer.Encrypt(sourceUserData)
	if err != nil {
		logger.DPanic("failed to encrypt the user data", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	jwtAccess, err := s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.ParamPurposeAccess)
	if err != nil {
		logger.DPanic("failed to generate an access token (JWT)", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	jwtRefresh, err := s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt.ParamPurposeRefresh)
	if err != nil {
		logger.DPanic("failed to generate an refresh token (JWT)", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.ServiceAccessTokenReturn{}, err
	}

	err = s.repository.UpdateSession(ctx, logger, model.RepoUpdateSessionParam{
		PublicSessionID: publicSessionID,
		RefreshToken:    refreshToken})
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

	decryptedJWTData, err := s.cryptographer.Decrypt(param.EncryptedUserData)
	if err != nil {
		logger.DPanic("failed to decrypt the user data", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = json.Unmarshal(decryptedJWTData, &user)
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

func (s *service) RequestUserPasswordReset(ctx context.Context, logger *zap.Logger, email string) error {

	remoteAddr, err := s.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return err
	}

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	confirmationKey := random.StringRand(16, 16, true)

	user, err := s.repository.RequestUserPasswordReset(ctx, logger, model.RepoRequestUserPasswordResetParam{
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

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: email},
		s.languageContent.Data.User.ResetPassword.Email.Request.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.ResetPassword.Email.Request.Body[s.languageContent.Language[user.Language]],
			remoteAddr, s.config.DomainAPI, confirmationKey),
		fmt.Sprintf(
			"<%s@%s>",
			confirmationKey,
			fmt.Sprintf("%s.%s", "reset-password", s.config.DomainAPI)))
	if err != nil {
		logger.DPanic("failed to add an email message to the queue", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (s *service) GetListActiveSessions(ctx context.Context, logger *zap.Logger, encryptedUserData []byte) ([]model.Session, error) {

	var user model.User

	requestID, requestIDKey, err := s.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return nil, err
	}

	decryptedJWTData, err := s.cryptographer.Decrypt(encryptedUserData)
	if err != nil {
		logger.DPanic("failed to decrypt the user data", zap.Error(err), zap.String(requestIDKey, requestID))
		return nil, err
	}

	err = json.Unmarshal(decryptedJWTData, &user)
	if err != nil {
		logger.DPanic("failed to unmarshal the decryptedJWTData value to struct", zap.Error(err),
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