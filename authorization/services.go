/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	cryptographer2 "github.com/dmalix/financelime-authorization/packages/cryptographer"
	email2 "github.com/dmalix/financelime-authorization/packages/email"
	jwt2 "github.com/dmalix/financelime-authorization/packages/jwt"
	"github.com/dmalix/financelime-authorization/utils/random"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"hash"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

func NewService(
	config ConfigService,
	languageContent LanguageContent,
	messageQueue chan email2.EmailMessage,
	message email2.Message,
	repository Repository,
	cryptographer cryptographer2.Cryptographer,
	jwt jwt2.Jwt) *service {
	return &service{
		config:          config,
		languageContent: languageContent,
		messageQueue:    messageQueue,
		message:         message,
		repository:      repository,
		cryptographer:   cryptographer,
		jwt:             jwt,
	}
}

func (s *service) generatePublicID(privateID int64) (string, error) {

	var (
		err             error
		hs              hash.Hash
		publicSessionID string
	)

	hs = sha256.New()
	_, err = hs.Write([]byte(
		strconv.FormatInt(privateID, 10) +
			random.StringRand(16, 16, false) +
			time.Now().String() +
			s.config.CryptoSalt))
	if err != nil {
		return publicSessionID, errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	publicSessionID = hex.EncodeToString(hs.Sum(nil))

	return publicSessionID, nil
}

func (s *service) signUp(param serviceSignUpParam) error {

	var (
		confirmationKey string
		err             error
	)

	confirmationKey = random.StringRand(16, 16, true)

	err = s.repository.createUser(repoCreateUserParam{
		email:              param.email,
		language:           param.language,
		inviteCode:         param.inviteCode,
		remoteAddr:         param.remoteAddr,
		confirmationKey:    confirmationKey,
		inviteCodeRequired: s.config.AuthInviteCodeRequired,
	})
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case domainErrorCodeBadParamEmail, domainErrorCodeBadParamInvite, domainErrorCodeBadParamLang,
			domainErrorCodeBadParamRemoteAddr, domainErrorCodeBadParamConfirmationKey:
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCodeBadParams,
				"one or more of the input parameters are invalid",
				err))
		case domainErrorCodeUserAlreadyExist:
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"a user with the email you specified already exists",
				err))
		case domainErrorCodeInviteNotFound:
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"the invite code does not exist or is expired",
				err))
		case domainErrorCodeInviteHasEnded:
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				domainErrorCode,
				"the limit for issuing this invite code has been exhausted",
				err))

		default:
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"a system error was returned",
				err))
		}
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: param.email},
		s.languageContent.Data.User.Signup.Email.Confirm.Subject[s.languageContent.Language[param.language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Signup.Email.Confirm.Body[s.languageContent.Language[param.language]],
			s.config.DomainAPI, confirmationKey),
		fmt.Sprintf(
			"<%s@%s>",
			confirmationKey,
			fmt.Sprintf("%s.%s", "sign-up", s.config.DomainAPI)))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to send message to the user",
			err))
	}

	return nil
}

func (s *service) confirmUserEmail(confirmationKey string) (string, error) {

	var (
		user user
		err  error
	)

	user, err = s.repository.confirmUserEmail(confirmationKey)

	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case domainErrorCodeBadParamConfirmationKey, domainErrorCodeConfirmationKeyNotFound,
			domainErrorCodeConfirmationKeyAlreadyConfirmed:
			return "",
				errors.New(fmt.Sprintf("%s:%s[%s]",
					domainErrorCodeBadConfirmationKey,
					"the confirmation key not valid",
					err))

		default:
			return "",
				errors.New(fmt.Sprintf("%s:%s[%s]",
					trace.GetCurrentPoint(),
					"a system error was returned",
					err))
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
		return "",
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to send message to the user",
				err))
	}

	confirmationMessage := s.languageContent.Data.User.Signup.Page.Text[s.languageContent.Language[user.Language]]

	return confirmationMessage, nil
}

func (s *service) createAccessToken(param serviceCreateAccessTokenParam) (serviceAccessTokenReturn, error) {

	var (
		err               error
		user              user
		sourceUserData    []byte
		encryptedUserData []byte
		publicSessionID   string
		accessJWT         string
		refreshJWT        string
	)

	user, err = s.repository.getUserByAuth(repoGetUserByAuthParam{
		email:    param.email,
		password: param.password,
	})
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case domainErrorCodeBadParamEmail, domainErrorCodeBadParamPassword, domainErrorCodeBadParamLang:
			return serviceAccessTokenReturn{},
				errors.New(fmt.Sprintf("%s:%s",
					domainErrorCodeBadParams,
					err))
		case domainErrorCodeUserNotFound:
			return serviceAccessTokenReturn{},
				errors.New(fmt.Sprintf("%s:%s",
					domainErrorCode,
					err))
		default:
			return serviceAccessTokenReturn{},
				errors.New(fmt.Sprintf("%s:%s[%s]",
					trace.GetCurrentPoint(),
					"A system error was returned",
					err))
		}
	}

	publicSessionID, err = s.generatePublicID(user.ID)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate the publicSessionID value",
				err))
	}

	sourceUserData, err = json.Marshal(user)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to marshal the user struct",
				err))
	}

	encryptedUserData, err = s.cryptographer.Encrypt(sourceUserData)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to encrypt the user data",
				err))
	}

	accessJWT, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt2.PropsPurposeAccess)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate an access token (JWT)",
				err))
	}

	refreshJWT, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt2.PropsPurposeRefresh)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate an refresh token (JWT)",
				err))
	}

	err = s.repository.saveSession(repoSaveSessionParam{
		userID:          user.ID,
		publicSessionID: publicSessionID,
		refreshToken:    refreshJWT,
		clientID:        param.clientID,
		remoteAddr:      param.remoteAddr,
		userAgent:       param.userAgent,
		device:          param.device,
	})
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to save the session",
				err))
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: param.email},
		s.languageContent.Data.User.Login.Email.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.Login.Email.Body[s.languageContent.Language[user.Language]],
			time.Now().UTC().String(),
			param.device.Platform,
			param.remoteAddr,
			s.config.DomainAPP),
		fmt.Sprintf(
			"<%s@%s>",
			param.remoteAddr,
			fmt.Sprintf("%s.%s", "get-access-token", s.config.DomainAPI)))
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to send message to the user",
				err))
	}

	return serviceAccessTokenReturn{
		publicSessionID: publicSessionID,
		accessJWT:       accessJWT,
		refreshJWT:      refreshJWT,
	}, nil
}

func (s *service) refreshAccessToken(param serviceRefreshAccessTokenParam) (serviceAccessTokenReturn, error) {

	var (
		err               error
		user              user
		sourceUserData    []byte
		encryptedUserData []byte
		jwtData           jwt2.JwtData
		publicSessionID   string
		jwtAccess         string
		jwtRefresh        string
	)

	jwtData, err = s.jwt.VerifyToken(param.refreshToken)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeBadRefreshToken,
				err))
	}

	user, err = s.repository.getUserByRefreshToken(param.refreshToken)
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case domainErrorCodeUserNotFound:
			return serviceAccessTokenReturn{},
				errors.New(fmt.Sprintf("%s:%s",
					domainErrorCode,
					err))
		default:
			return serviceAccessTokenReturn{},
				errors.New(fmt.Sprintf("%s:%s[%s]",
					trace.GetCurrentPoint(),
					"A system error was returned",
					err))
		}
	}

	publicSessionID = jwtData.Payload.PublicSessionID

	sourceUserData, err = json.Marshal(user)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to marshal the user struct",
				err))
	}

	encryptedUserData, err = s.cryptographer.Encrypt(sourceUserData)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to encrypt the user data",
				err))
	}

	jwtAccess, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt2.PropsPurposeAccess)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate an access token (JWT)",
				err))
	}

	jwtRefresh, err = s.jwt.GenerateToken(publicSessionID, encryptedUserData, jwt2.PropsPurposeRefresh)
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to generate an refresh token (JWT)",
				err))
	}

	err = s.repository.updateSession(repoUpdateSessionParam{
		publicSessionID: publicSessionID,
		refreshToken:    param.refreshToken,
		remoteAddr:      param.remoteAddr,
	})
	if err != nil {
		return serviceAccessTokenReturn{},
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to save the session",
				err))
	}

	return serviceAccessTokenReturn{
		publicSessionID: publicSessionID,
		accessJWT:       jwtAccess,
		refreshJWT:      jwtRefresh}, nil
}

/*
	Revoke Access Token
		----------------
		Return:
			err             error  - system error
*/
func (s *service) revokeRefreshToken(param serviceRevokeRefreshTokenParam) error {

	var (
		err               error
		user              user
		decryptedUserData []byte
	)

	decryptedUserData, err = s.cryptographer.Decrypt(param.encryptedUserData)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to decrypt the user data",
			err))
	}

	err = json.Unmarshal(decryptedUserData, &user)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to unmarshal the decryptedUserData value to struct [%s]",
			err))
	}

	err = s.repository.deleteSession(repoDeleteSessionParam{
		userID:          user.ID,
		publicSessionID: param.publicSessionID,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to detele the session [%s]",
			err))
	}

	return nil
}

func (s *service) requestUserPasswordReset(param serviceRequestUserPasswordResetParam) error {

	var (
		confirmationKey string
		err             error
		user            user
	)

	confirmationKey = random.StringRand(16, 16, true)

	user, err = s.repository.requestUserPasswordReset(repoRequestUserPasswordResetParam{
		email:           param.email,
		remoteAddr:      param.remoteAddr,
		confirmationKey: confirmationKey})
	if err != nil {
		domainErrorCode := strings.Split(err.Error(), ":")[0]
		switch domainErrorCode {
		case domainErrorCodeBadParamEmail, domainErrorCodeBadParamRemoteAddr, domainErrorCodeBadConfirmationKey:
			return errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeBadParamEmail,
				err))
		case domainErrorCodeUserNotFound: // a user with the email specified not found
			return errors.New(fmt.Sprintf("%s:%s",
				domainErrorCode,
				err))
		default:
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"a system error was returned",
				err))
		}
	}

	err = s.message.AddEmailMessageToQueue(
		s.messageQueue,
		mail.Address{Address: param.email},
		s.languageContent.Data.User.ResetPassword.Email.Request.Subject[s.languageContent.Language[user.Language]],
		fmt.Sprintf(
			s.languageContent.Data.User.ResetPassword.Email.Request.Body[s.languageContent.Language[user.Language]],
			param.remoteAddr, s.config.DomainAPI, confirmationKey),
		fmt.Sprintf(
			"<%s@%s>",
			confirmationKey,
			fmt.Sprintf("%s.%s", "reset-password", s.config.DomainAPI)))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to send message to the user",
			err))
	}

	return nil
}

func (s *service) getListActiveSessions(encryptedUserData []byte) ([]session, error) {

	var (
		err               error
		sessions          []session
		decryptedUserData []byte
		user              user
	)

	decryptedUserData, err = s.cryptographer.Decrypt(encryptedUserData)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to decrypt the user data",
				err))
	}

	err = json.Unmarshal(decryptedUserData, &user)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"Failed to unmarshal the decryptedUserData value to struct [%s]",
				err))
	}

	sessions, err = s.repository.getListActiveSessions(user.ID)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				trace.GetCurrentPoint(),
				"a system error was returned",
				err))
	}
	return sessions, nil
}
