/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/dmalix/financelime-authorization/app/authorization"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"github.com/dmalix/financelime-authorization/packages/generator"
	"github.com/dmalix/financelime-authorization/packages/middleware"
	"go.uber.org/zap"
	"html"
	"regexp"
	"strconv"
	"strings"
)

type repository struct {
	config        model.ConfigRepository
	contextGetter middleware.ContextGetter
	dbAuthMain    *sql.DB
	dbAuthRead    *sql.DB
	dbBlade       *sql.DB
}

func NewRepository(
	config model.ConfigRepository,
	contextGetter middleware.ContextGetter,
	dbAuthMain,
	dbAuthRead,
	dbBlade *sql.DB) *repository {
	return &repository{
		config:        config,
		contextGetter: contextGetter,
		dbAuthMain:    dbAuthMain,
		dbAuthRead:    dbAuthRead,
		dbBlade:       dbBlade,
	}
}

func (r *repository) SignUpStep1(ctx context.Context, logger *zap.Logger, param model.RepoSignUpParam) error {

	var (
		userID                   int64
		inviteCodeRecord         model.InviteCodeRecord
		inviteCodeReservedID     int64
		amountInviteCodeIssued   int
		amountInviteCodeReserved int
		inviteCodesIsRunOut      bool
		confirmationID           int64
	)

	remoteAddr, _, err := r.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return err
	}

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	// Check parameters

	param.Email = html.EscapeString(param.Email)
	if len(param.Email) <= 2 || len(param.Email) > 255 {
		logger.Error("the param is not valid", zap.String("email", param.Email),
			zap.String(requestIDKey, requestID))
		return authorization.ErrorBadParamEmail
	}

	param.InviteCode = html.EscapeString(param.InviteCode)
	paramValueRegexp := regexp.MustCompile(`^[0-9a-zA-Z_-]{3,16}$`)
	if !paramValueRegexp.MatchString(param.InviteCode) {
		if param.InviteCodeRequired {
			logger.Error("the param is not valid", zap.String("inviteCode", param.InviteCode),
				zap.String(requestIDKey, requestID))
			return authorization.ErrorBadParamInvite
		}
	}

	param.Language = html.EscapeString(param.Language)
	paramValueRegexp = regexp.MustCompile(`^[ru|en]{2}$`)
	if !paramValueRegexp.MatchString(param.Language) {
		logger.Error("the param is not valid", zap.String("language", param.Language),
			zap.String(requestIDKey, requestID))
		return authorization.ErrorBadParamLang
	}

	param.ConfirmationKey = html.EscapeString(param.ConfirmationKey)
	paramValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(param.ConfirmationKey) {
		logger.Error("the confirmationKey param is not valid", zap.String("confirmationKey", param.ConfirmationKey),
			zap.String(requestIDKey, requestID))
		return authorization.ErrorBadParamConfirmationKey
	}

	// Begin the transaction

	dbTransactionAuthMain, err := r.dbAuthMain.Begin()
	if err != nil {
		logger.DPanic("failed to begin AuthMain DB transaction", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}
	defer func(dbTransactionAuthMain *sql.Tx) {
		logger.Debug("the AuthMain DB transaction was rollback", zap.String(requestIDKey, requestID))
		err := dbTransactionAuthMain.Rollback()
		if err != nil && err.Error() != messageTransactionHasAlreadyBeenCommittedOrRolledBack {
			logger.DPanic("failed to rollback AuthMain DB transaction", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(dbTransactionAuthMain)

	dbTransactionBlade, err := r.dbBlade.Begin()
	if err != nil {
		logger.DPanic("failed to begin Blade DB transaction", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}
	defer func(dbTransactionBlade *sql.Tx) {
		logger.Debug("the Blade DB transaction was rollback", zap.String(requestIDKey, requestID))
		err := dbTransactionBlade.Rollback()
		if err != nil && err.Error() != messageTransactionHasAlreadyBeenCommittedOrRolledBack {
			logger.DPanic("failed to rollback Blade DB transaction", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(dbTransactionBlade)

	// Lock tables

	_, err = dbTransactionAuthMain.Exec("/* postgreSQL query */\n" +
		"LOCK TABLE \"user\",\n" +
		"invite_code,\n" +
		"invite_code_issued IN SHARE ROW EXCLUSIVE MODE\n")
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	_, err = dbTransactionBlade.Exec("/* postgreSQL query */\n" +
		"LOCK TABLE confirmation_create_new_user,\n" +
		"invite_code_reserved IN SHARE ROW EXCLUSIVE MODE\n")
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	if len(param.InviteCode) > 0 {

		// Check if this Invite Code exists

		loadInviteCode, err := dbTransactionAuthMain.Query("/* postgreSQL query */\n"+
			"SELECT\n"+
			"    invite_code.\"id\",\n"+
			"    invite_code.number_limit,\n"+
			"    \"user\".\"id\" AS user_id\n"+
			"FROM\n"+
			"    invite_code\n"+
			"INNER JOIN \"user\" ON\n"+
			"    invite_code.user_id = \"user\".\"id\"\n"+
			"WHERE\n"+
			"    invite_code.\"value\" = $1\n"+
			"    AND \"user\".deleted_at IS NULL\n"+
			"    AND invite_code.deleted_at IS NULL\n"+
			"    AND invite_code.expires_at > NOW( )\n"+
			"LIMIT 1\n",
			param.InviteCode)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return err
		}
		defer func(loadInviteCode *sql.Rows) {
			if err := loadInviteCode.Close(); err != nil {
				logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
			}
		}(loadInviteCode)

		for loadInviteCode.Next() {
			err := loadInviteCode.Scan(&inviteCodeRecord.Id, &inviteCodeRecord.LimitAmount, &inviteCodeRecord.UserID)
			if err != nil {
				logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
				return err
			}
		}

		if inviteCodeRecord.Id == 0 {
			logger.Error("the invite code does not exist or is expired", zap.String(requestIDKey, requestID))
			return authorization.ErrorInviteNotFound
		}

		// Check the limit for this invite code, including the reservation

		loadAmountInviteCodeIssued, err := dbTransactionAuthMain.Query("/* postgreSQL query */\n"+
			"SELECT\n"+
			"    COUNT (\n"+
			"        invite_code_issued.\"id\"\n"+
			"    )\n"+
			"FROM\n"+
			"    invite_code\n"+
			"INNER JOIN \"user\" ON\n"+
			"    invite_code.user_id = \"user\".\"id\"\n"+
			"INNER JOIN invite_code_issued ON\n"+
			"    invite_code.\"id\" = invite_code_issued.invite_code_id\n"+
			"WHERE\n"+
			"    invite_code.\"id\" = $1\n"+
			"    AND \"user\".deleted_at IS NULL\n"+
			"    AND invite_code_issued.deleted_at IS NULL\n"+
			"    AND invite_code.deleted_at IS NULL\n"+
			"    AND invite_code.expires_at > NOW( )\n",
			inviteCodeRecord.Id)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return err
		}
		defer func(loadAmountInviteCodeIssued *sql.Rows) {
			if err := loadAmountInviteCodeIssued.Close(); err != nil {
				logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
			}
		}(loadAmountInviteCodeIssued)

		for loadAmountInviteCodeIssued.Next() {
			err = loadAmountInviteCodeIssued.Scan(&amountInviteCodeIssued)
			if err != nil {
				logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
				return err
			}
		}

		loadAmountInviteCodeReserved, err := dbTransactionBlade.Query("/* postgreSQL query */\n"+
			"SELECT\n"+
			"    count(invite_code_reserved.\"id\")\n"+
			"FROM\n"+
			"    invite_code_reserved\n"+
			"INNER JOIN confirmation_create_new_user ON\n"+
			"    invite_code_reserved.email = confirmation_create_new_user.email\n"+
			"WHERE\n"+
			"    invite_code_reserved.invite_code_id = $1\n"+
			"    AND invite_code_reserved.deleted_at IS NULL\n"+
			"    AND confirmation_create_new_user.deleted_at IS NULL\n"+
			"    AND confirmation_create_new_user.expires_at > NOW()\n",
			inviteCodeRecord.Id)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return err
		}
		defer func(loadAmountInviteCodeReserved *sql.Rows) {
			if err := loadAmountInviteCodeReserved.Close(); err != nil {
				logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
			}
		}(loadAmountInviteCodeReserved)

		for loadAmountInviteCodeReserved.Next() {
			err = loadAmountInviteCodeReserved.Scan(&amountInviteCodeReserved)
			if err != nil {
				logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
				return err
			}
		}

		if (amountInviteCodeIssued + amountInviteCodeReserved) >= inviteCodeRecord.LimitAmount {
			inviteCodesIsRunOut = true

			if param.InviteCodeRequired {
				logger.Error("the invite code has ended", zap.String(requestIDKey, requestID))
				return authorization.ErrorInviteHasEnded
			}
		}
	}

	// Check if an user exists with this email, including new users pending confirmation

	loadVerifiedUser, err := dbTransactionAuthMain.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"user\".id\n"+
		"FROM\n"+
		"    \"user\"\n"+
		"WHERE\n"+
		"    \"user\".email = $1\n"+
		"    AND \"user\".deleted_at IS NULL\n"+
		"LIMIT 1\n", param.Email)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}
	defer func(loadVerifiedUser *sql.Rows) {
		if err := loadVerifiedUser.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadVerifiedUser)

	for loadVerifiedUser.Next() {
		err = loadVerifiedUser.Scan(&userID)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return err
		}
	}

	loadPendingUser, err := dbTransactionBlade.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    confirmation_create_new_user.\"id\"\n"+
		"FROM\n"+
		"    confirmation_create_new_user\n"+
		"WHERE\n"+
		"    confirmation_create_new_user.email = $1\n"+
		"    AND confirmation_create_new_user.deleted_at IS NULL\n"+
		"    AND confirmation_create_new_user.expires_at > NOW( )\n"+
		"LIMIT 1\n",
		param.Email)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}
	defer func(loadPendingUser *sql.Rows) {
		if err := loadPendingUser.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadPendingUser)

	for loadPendingUser.Next() {
		err = loadPendingUser.Scan(&confirmationID)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return err
		}
	}

	if userID != 0 || confirmationID != 0 {
		logger.Error("a user with the same email address already exists", zap.String(requestIDKey, requestID))
		return authorization.ErrorUserAlreadyExist
	}

	// Add a new record for the user pending confirmation
	// TODO Move the 1440 value of interval to config
	err = dbTransactionBlade.QueryRow("/* postgreSQL query */\n"+
		"INSERT\n"+
		"    INTO\n"+
		"    confirmation_create_new_user (\n"+
		"        created_at,\n"+
		"        email,\n"+
		"        \"language\",\n"+
		"        confirmation_key,\n"+
		"        remote_addr,\n"+
		"        expires_at\n"+
		"    )\n"+
		"VALUES (\n"+
		"    NOW( ),\n"+
		"    $1,\n"+
		"    $2,\n"+
		"    $3,\n"+
		"    $4,\n"+
		"    NOW( ) + INTERVAL '1440 minute'\n"+
		") RETURNING \"id\"\n",
		param.Email, param.Language, param.ConfirmationKey, remoteAddr).Scan(&confirmationID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	if len(param.InviteCode) > 0 && !inviteCodesIsRunOut {

		// Add a new record for reserve the invite code for the user

		err = dbTransactionBlade.QueryRow("/* postgreSQL query */\n"+
			"INSERT\n"+
			"    INTO\n"+
			"    invite_code_reserved (\n"+
			"        created_at,\n"+
			"        invite_code_id,\n"+
			"        email,\n"+
			"        confirmation_id\n"+
			"    )\n"+
			"VALUES (\n"+
			"    NOW( ),\n"+
			"    $1,\n"+
			"    $2,\n"+
			"    $3\n"+
			") RETURNING \"id\"\n",
			inviteCodeRecord.Id, param.Email, confirmationID).Scan(&inviteCodeReservedID)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return err
		}
	}

	// Transactions Commit

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		logger.DPanic("failed to commit to the AuthMain DB", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		logger.DPanic("failed to commit to the Blade DB", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (r *repository) SignUpStep2(ctx context.Context, logger *zap.Logger, confirmationKey string) (model.User, error) {

	var (
		confirmationID     int64
		inviteCodeIssuedID int64
		inviteCodeID       int64
		userID             int64
		user               model.User
	)

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.User{}, err
	}

	propsValueRegexp := regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !propsValueRegexp.MatchString(confirmationKey) {
		logger.Error("the param is not valid", zap.String("confirmationKey", confirmationKey))
		return model.User{}, authorization.ErrorBadParamConfirmationKey
	}

	// Check the confirmationKey in Database

	loadConfirmationKeyInfo, err := r.dbBlade.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    confirmation_create_new_user.\"id\",\n"+
		"    confirmation_create_new_user.email,\n"+
		"    confirmation_create_new_user.\"language\"\n"+
		"FROM\n"+
		"    confirmation_create_new_user\n"+
		"WHERE\n"+
		"    confirmation_create_new_user.confirmation_key = $1\n"+
		"    AND confirmation_create_new_user.deleted_at IS NULL\n"+
		"    AND confirmation_create_new_user.expires_at > NOW( )\n"+
		"ORDER BY\n"+
		"    confirmation_create_new_user.\"id\" DESC\n"+
		"LIMIT 1\n",
		confirmationKey)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID),
			zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(loadConfirmationKeyInfo *sql.Rows) {
		if err := loadConfirmationKeyInfo.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadConfirmationKeyInfo)

	for loadConfirmationKeyInfo.Next() {
		err = loadConfirmationKeyInfo.Scan(&confirmationID, &user.Email, &user.Language)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if confirmationID == 0 {
		return model.User{}, authorization.ErrorConfirmationKeyNotFound
	}

	// Begin the transaction

	dbTransactionAuthMain, err := r.dbAuthMain.Begin()
	if err != nil {
		logger.DPanic("failed to begin AuthMain DB transaction", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(dbTransactionAuthMain *sql.Tx) {
		logger.Debug("the AuthMain DB transaction was rollback", zap.String(requestIDKey, requestID))
		err := dbTransactionAuthMain.Rollback()
		if err != nil && err.Error() != messageTransactionHasAlreadyBeenCommittedOrRolledBack {
			logger.DPanic("failed to rollback AuthMain DB transaction", zap.Error(err),
				zap.String(requestIDKey, requestID))
		}
	}(dbTransactionAuthMain)

	dbTransactionBlade, err := r.dbBlade.Begin()
	if err != nil {
		logger.DPanic("failed to begin Blade DB transaction", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(dbTransactionBlade *sql.Tx) {
		logger.Debug("the Blade DB transaction was rollback", zap.String(requestIDKey, requestID))
		err := dbTransactionBlade.Rollback()
		if err != nil && err.Error() != messageTransactionHasAlreadyBeenCommittedOrRolledBack {
			logger.DPanic("failed to rollback Blade DB transaction", zap.Error(err),
				zap.String(requestIDKey, requestID))
		}
	}(dbTransactionBlade)

	// Lock tables

	_, err = dbTransactionAuthMain.Exec("/* postgreSQL query */\n" +
		"LOCK TABLE \"user\",\n" +
		"invite_code,\n" +
		"invite_code_issued IN SHARE ROW EXCLUSIVE MODE\n")
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	_, err = dbTransactionBlade.Exec("/* postgreSQL query */\n" +
		"LOCK TABLE confirmation_create_new_user,\n" +
		"invite_code_reserved IN SHARE ROW EXCLUSIVE MODE\n")
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	// Check if an user exists with this email address

	loadUserID, err := dbTransactionAuthMain.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"user\".\"id\"\n"+
		"FROM\n"+
		"    \"user\"\n"+
		"WHERE\n"+
		"    \"user\".email = $1\n"+
		"    AND \"user\".deleted_at IS NULL\n"+
		"LIMIT 1\n",
		user.Email)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(loadUserID *sql.Rows) {
		if err := loadUserID.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadUserID)

	for loadUserID.Next() {
		err = loadUserID.Scan(&userID)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if userID != 0 {
		return model.User{}, authorization.ErrorConfirmationKeyAlreadyConfirmed
	}

	// Updating the confirmation key status to "Deleted"

	deleteConfirmationKey, err := dbTransactionBlade.Prepare("/* postgreSQL query */\n" +
		"UPDATE\n" +
		"    confirmation_create_new_user\n" +
		"SET\n" +
		"    deleted_at = NOW()\n" +
		"WHERE\n" +
		"    confirmation_create_new_user.\"id\" = $1\n")
	if err != nil {
		logger.DPanic("failed to prepare the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(deleteConfirmationKey *sql.Stmt) {
		if err := deleteConfirmationKey.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(deleteConfirmationKey)

	_, err = deleteConfirmationKey.Exec(confirmationID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	// Check if there was a reserve for an invite code

	loadInviteCode, err := dbTransactionBlade.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    invite_code_reserved.invite_code_id\n"+
		"FROM\n"+
		"    invite_code_reserved\n"+
		"WHERE\n"+
		"    invite_code_reserved.confirmation_id = $1\n"+
		"    AND invite_code_reserved.deleted_at IS NULL\n"+
		"LIMIT 1\n",
		confirmationID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(loadInviteCode *sql.Rows) {
		if err := loadInviteCode.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadInviteCode)

	for loadInviteCode.Next() {
		err = loadInviteCode.Scan(&inviteCodeID)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if inviteCodeID > 0 {

		// Updating the invite code reserve to "Deleted"

		deleteInviteCode, err := dbTransactionBlade.Prepare("/* postgreSQL query */\n" +
			"UPDATE\n" +
			"    invite_code_reserved\n" +
			"SET\n" +
			"    deleted_at = NOW()\n" +
			"WHERE\n" +
			"    invite_code_reserved.confirmation_id = $1\n")
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
		defer func(deleteInviteCode *sql.Stmt) {
			if err := deleteInviteCode.Close(); err != nil {
				logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
			}
		}(deleteInviteCode)

		_, err = deleteInviteCode.Exec(confirmationID)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	// Create the new verified user in the Auth DB

	user.Password = generator.StringRand(16, 16, false)
	hs := sha256.New()
	_, err = hs.Write([]byte(user.Password + r.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate password for new user", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	hashedPassword := hex.EncodeToString(hs.Sum(nil))

	err = dbTransactionAuthMain.QueryRow("/* postgreSQL query */\n"+
		"INSERT\n"+
		"    INTO\n"+
		"    \"user\" (\n"+
		"        created_at,\n"+
		"        email,\n"+
		"        \"language\",\n"+
		"        \"password\"\n"+
		"    )\n"+
		"VALUES (\n"+
		"    NOW( ),\n"+
		"    $1,\n"+
		"    $2,\n"+
		"    $3\n"+
		") RETURNING \"id\"\n",
		user.Email, user.Language, hashedPassword).
		Scan(&user.ID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	if inviteCodeID > 0 {

		// Linking the invite code to the user

		err = dbTransactionAuthMain.QueryRow("/* postgreSQL query */\n"+
			"INSERT\n"+
			"    INTO\n"+
			"    invite_code_issued (\n"+
			"        created_at,\n"+
			"        invite_code_id,\n"+
			"        user_id\n"+
			"    )\n"+
			"VALUES (\n"+
			"    NOW( ),\n"+
			"    $1,\n"+
			"    $2\n"+
			") RETURNING \"id\"\n",
			inviteCodeID, user.ID).
			Scan(&inviteCodeIssuedID)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	// Transactions Commit

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		logger.DPanic("failed to commit to the AuthMain DB", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		logger.DPanic("failed to commit to the Blade DB", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	return user, nil
}

func (r *repository) GetUserByAuth(ctx context.Context, logger *zap.Logger, param model.RepoGetUserByAuthParam) (model.User, error) {

	var user model.User

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.User{}, err
	}

	// Check parameters

	if len(param.Email) <= 2 || len(param.Email) > 255 {
		logger.Error("the param is not valid", zap.String("email", param.Email),
			zap.String(requestIDKey, requestID))
		return model.User{}, authorization.ErrorBadParamEmail
	}

	if len(param.Password) <= 2 || len(param.Password) > 255 {
		logger.Error("the param is not valid", zap.String("password", param.Password),
			zap.String(requestIDKey, requestID))
		return model.User{}, authorization.ErrorBadParamPassword
	}

	// Get User

	hs := sha256.New()
	_, err = hs.Write([]byte(param.Password + r.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate hash", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	hashedPassword := hex.EncodeToString(hs.Sum(nil))

	loadUser, err := r.dbAuthRead.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"user\".\"id\",\n"+
		"    \"user\".email,\n"+
		"    \"user\".\"language\"\n"+
		"FROM\n"+
		"    \"user\"\n"+
		"WHERE\n"+
		"    \"user\".email = $1\n"+
		"    AND \"user\".PASSWORD = $2\n"+
		"    AND \"user\".deleted_at IS NULL\n"+
		"LIMIT 1\n", param.Email, hashedPassword)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(loadUser *sql.Rows) {
		if err := loadUser.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadUser)

	for loadUser.Next() {
		err = loadUser.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if user.ID == 0 {
		logger.Error("user not found", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, authorization.ErrorUserNotFound
	}

	return user, nil
}

func (r *repository) GetUserByRefreshToken(ctx context.Context, logger *zap.Logger, refreshToken string) (model.User, error) {

	var user model.User

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.User{}, err
	}

	hs := sha256.New()
	_, err = hs.Write([]byte(
		refreshToken +
			r.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate hash", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	hashedRefreshToken := hex.EncodeToString(hs.Sum(nil))

	loadUser, err := r.dbAuthRead.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"user\".\"id\",\n"+
		"    \"user\".email,\n"+
		"    \"user\".\"language\"\n"+
		"FROM\n"+
		"    \"session\"\n"+
		"INNER JOIN \"user\" ON\n"+
		"    \"session\".user_id = \"user\".\"id\"\n"+
		"WHERE\n"+
		"    \"session\".hashed_refresh_token = $1\n"+
		"    AND \"session\".deleted_at IS NULL\n"+
		"    AND \"user\".deleted_at IS NULL\n"+
		"LIMIT 1\n", hashedRefreshToken)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return user, err
	}
	defer func(loadUser *sql.Rows) {
		if err := loadUser.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadUser)

	for loadUser.Next() {
		err = loadUser.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return user, err
		}
	}

	if user.ID == 0 {
		logger.DPanic("user not found", zap.Error(err), zap.String("hashedRefreshToken", hashedRefreshToken))
		return user, authorization.ErrorUserNotFound
	}

	return user, nil
}

func (r *repository) SaveSession(ctx context.Context, logger *zap.Logger, param model.RepoSaveSessionParam) error {

	var sessionID int64
	var deviceID int64

	remoteAddr, _, err := r.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return err
	}

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	hs := sha256.New()
	_, err = hs.Write([]byte(
		param.RefreshToken +
			r.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate hash for refresh token", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	hashedRefreshToken := hex.EncodeToString(hs.Sum(nil))

	err = r.dbAuthMain.QueryRow("/* postgreSQL query */\n"+
		"INSERT\n"+
		"    INTO\n"+
		"    \"session\" (\n"+
		"        created_at,\n"+
		"        user_id,\n"+
		"        client_id,\n"+
		"        remote_addr,\n"+
		"        public_id,\n"+
		"        hashed_refresh_token\n"+
		"    )\n"+
		"VALUES (\n"+
		"    NOW(),\n"+
		"    $1,\n"+
		"    $2,\n"+
		"    $3,\n"+
		"    $4,\n"+
		"    $5\n"+
		") RETURNING \"id\"\n",
		param.UserID,
		param.ClientID,
		remoteAddr,
		param.PublicSessionID,
		hashedRefreshToken).
		Scan(&sessionID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	err = r.dbAuthMain.QueryRow("/* postgreSQL query */\n"+
		"INSERT\n"+
		"    INTO\n"+
		"    device (\n"+
		"        created_at,\n"+
		"        session_id,\n"+
		"        platform,\n"+
		"        height,\n"+
		"        width,\n"+
		"        \"language\",\n"+
		"        timezone,\n"+
		"        user_agent\n"+
		"    )\n"+
		"VALUES (\n"+
		"    NOW(),\n"+
		"    $1,\n"+
		"    $2,\n"+
		"    $3,\n"+
		"    $4,\n"+
		"    $5,\n"+
		"    $6,\n"+
		"    $7\n"+
		") RETURNING \"id\"\n",
		sessionID,
		html.EscapeString(param.Device.Platform),
		param.Device.Height,
		param.Device.Width,
		html.EscapeString(param.Device.Language),
		html.EscapeString(param.Device.Timezone),
		html.EscapeString(param.UserAgent)).
		Scan(&deviceID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (r *repository) UpdateSession(ctx context.Context, logger *zap.Logger, param model.RepoUpdateSessionParam) error {

	var result string

	remoteAddr, _, err := r.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return err
	}

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	hs := sha256.New()
	_, err = hs.Write([]byte(
		param.RefreshToken +
			r.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate hash", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	hashedRefreshToken := hex.EncodeToString(hs.Sum(nil))

	err = r.dbAuthMain.QueryRow("/* postgreSQL query */\n"+
		"UPDATE\n"+
		"    \"session\"\n"+
		"SET\n"+
		"    updated_at = NOW( ),\n"+
		"    remote_addr = $1,\n"+
		"    hashed_refresh_token = $2\n"+
		"WHERE\n"+
		"    \"session\".public_id = $3\n"+
		"    AND \"session\".hashed_refresh_token = $4\n"+
		"    AND \"session\".deleted_at IS NULL RETURNING \"session\".public_id\n",
		remoteAddr,
		hashedRefreshToken,
		param.PublicSessionID,
		hashedRefreshToken).Scan(&result)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	if len(result) == 0 {
		logger.Error("the case (the session + hashedRefreshToken) does not exist",
			zap.String("publicSessionID", param.PublicSessionID), zap.String("hashedRefreshToken", hashedRefreshToken),
			zap.String(requestIDKey, requestID))
		return authorization.ErrorSessionNotFound
	}

	return nil
}

func (r *repository) GetListActiveSessions(ctx context.Context, logger *zap.Logger, userID int64) ([]model.Session, error) {

	var session model.Session
	var sessions []model.Session

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return nil, err
	}

	loadActiveSessionsList, err := r.dbAuthRead.Query(strings.Replace("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"session\".public_id,\n"+
		"    device.platform,\n"+
		"    CASE\n"+
		"        WHEN (\n"+
		"            \"session\".updated_at IS NULL\n"+
		"        ) THEN \"session\".created_at\n"+
		"        ELSE \"session\".updated_at\n"+
		"    END AS updated_at\n"+
		"FROM\n"+
		"    \"session\"\n"+
		"INNER JOIN device ON\n"+
		"    \"session\".\"id\" = device.session_id\n"+
		"WHERE\n"+
		"    (\n"+
		"        (\n"+
		"            \"session\".updated_at IS NULL\n"+
		"                AND \"session\".created_at > NOW( ) - INTERVAL '$LIFETIME SECOND'\n"+
		"        )\n"+
		"            OR (\n"+
		"                \"session\".updated_at IS NOT NULL\n"+
		"                    AND \"session\".updated_at > NOW( ) - INTERVAL '$LIFETIME SECOND'\n"+
		"            )\n"+
		"    )\n"+
		"    AND \"session\".user_id = $1\n"+
		"    AND \"session\".deleted_at IS NULL\n",
		"$LIFETIME", strconv.Itoa(r.config.JwtRefreshTokenLifetime), 2),
		userID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return nil, err
	}
	defer func(loadActiveSessionsList *sql.Rows) {
		if err := loadActiveSessionsList.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadActiveSessionsList)

	for loadActiveSessionsList.Next() {
		err = loadActiveSessionsList.Scan(
			&session.PublicSessionID,
			&session.Platform,
			&session.UpdatedAt,
		)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *repository) DeleteSession(ctx context.Context, logger *zap.Logger, param model.RepoDeleteSessionParam) error {

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return err
	}

	_, err = r.dbAuthMain.Exec("/* postgreSQL query */\n"+
		"UPDATE\n"+
		"    \"session\"\n"+
		"SET\n"+
		"    deleted_at = NOW()\n"+
		"WHERE\n"+
		"    \"session\".user_id = $1\n"+
		"    AND \"session\".public_id = $2\n",
		param.UserID,
		param.PublicSessionID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return err
	}

	return nil
}

func (r *repository) ResetUserPasswordStep1(ctx context.Context, logger *zap.Logger, param model.RepoResetUserPasswordParam) (model.User, error) {

	var user model.User
	var confirmationID string

	remoteAddr, _, err := r.contextGetter.GetRemoteAddr(ctx)
	if err != nil {
		logger.DPanic("failed to get remoteAddr", zap.Error(err))
		return model.User{}, err
	}

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.User{}, err
	}

	// Check props

	param.Email = html.EscapeString(param.Email)
	if len(param.Email) <= 2 || len(param.Email) > 255 {
		logger.Error("the param is not valid", zap.String("email", param.Email),
			zap.String(requestIDKey, requestID))
		return user, authorization.ErrorBadParamEmail
	}

	param.ConfirmationKey = html.EscapeString(param.ConfirmationKey)
	paramValueRegexp := regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(param.ConfirmationKey) {
		logger.Error("the param is not valid", zap.String("confirmationKey", param.ConfirmationKey),
			zap.String(requestIDKey, requestID))
		return user, authorization.ErrorBadParamEmail
	}

	// Check if the user exists

	loadUser, err := r.dbAuthRead.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"user\".\"id\",\n"+
		"    \"user\".email,\n"+
		"    \"user\".\"language\"\n"+
		"FROM\n"+
		"    \"user\"\n"+
		"WHERE\n"+
		"    \"user\".email = $1\n"+
		"    AND \"user\".deleted_at IS NULL\n"+
		"LIMIT 1\n", param.Email)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return user, err
	}
	defer func(loadUser *sql.Rows) {
		if err := loadUser.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadUser)

	for loadUser.Next() {
		err = loadUser.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if user.ID == 0 {
		logger.Error("the user not found", zap.String(requestIDKey, requestID))
		return model.User{}, authorization.ErrorUserNotFound

	}

	// Add the new record about reset password
	// TODO Move the 15 value of interval to config
	err = r.dbBlade.QueryRow("/* postgreSQL query */\n"+
		"INSERT\n"+
		"    INTO\n"+
		"    confirmation_reset_password (\n"+
		"        created_at,\n"+
		"        email,\n"+
		"        \"language\",\n"+
		"        confirmation_key,\n"+
		"        remote_addr,\n"+
		"        expires_at\n"+
		"    )\n"+
		"VALUES (\n"+
		"    NOW(),\n"+
		"    $1,\n"+
		"    $2,\n"+
		"    $3,\n"+
		"    $4,\n"+
		"    NOW() + INTERVAL '15 minute'\n"+
		") RETURNING \"id\"\n",
		param.Email,
		user.Language,
		param.ConfirmationKey,
		remoteAddr).
		Scan(&confirmationID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	return user, nil
}

func (r *repository) ResetUserPasswordStep2(ctx context.Context, logger *zap.Logger, confirmationKey string) (model.User, error) {

	var (
		confirmationKeyID int64
		userID            int64
		user              model.User
	)

	requestID, requestIDKey, err := r.contextGetter.GetRequestID(ctx)
	if err != nil {
		logger.DPanic("failed to get requestID", zap.Error(err))
		return model.User{}, err
	}

	propsValueRegexp := regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !propsValueRegexp.MatchString(confirmationKey) {
		logger.Error("the param is not valid", zap.String("confirmationKey", confirmationKey))
		return model.User{}, authorization.ErrorBadParamConfirmationKey
	}

	// Check the confirmationKey in Database

	loadConfirmationKeyInfo, err := r.dbBlade.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    confirmation_reset_password.\"id\",\n"+
		"    confirmation_reset_password.email,\n"+
		"    confirmation_reset_password.\"language\"\n"+
		"FROM\n"+
		"    confirmation_reset_password\n"+
		"WHERE\n"+
		"    confirmation_reset_password.confirmation_key = $1\n"+
		"    AND confirmation_reset_password.deleted_at IS NULL\n"+
		"    AND confirmation_reset_password.expires_at > NOW()\n"+
		"ORDER BY\n"+
		"    confirmation_reset_password.\"id\" DESC\n"+
		"LIMIT 1\n",
		confirmationKey)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID),
			zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(loadConfirmationKeyInfo *sql.Rows) {
		if err := loadConfirmationKeyInfo.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadConfirmationKeyInfo)

	for loadConfirmationKeyInfo.Next() {
		err = loadConfirmationKeyInfo.Scan(&confirmationKeyID, &user.Email, &user.Language)
		if err != nil {
			logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if confirmationKeyID == 0 {
		return model.User{}, authorization.ErrorConfirmationKeyNotFound
	}

	// Begin the transaction

	dbTransactionAuthMain, err := r.dbAuthMain.Begin()
	if err != nil {
		logger.DPanic("failed to begin AuthMain DB transaction", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(dbTransactionAuthMain *sql.Tx) {
		logger.Debug("the AuthMain DB transaction was rollback", zap.String(requestIDKey, requestID))
		err := dbTransactionAuthMain.Rollback()
		if err != nil && err.Error() != messageTransactionHasAlreadyBeenCommittedOrRolledBack {
			logger.DPanic("failed to rollback AuthMain DB transaction", zap.Error(err),
				zap.String(requestIDKey, requestID))
		}
	}(dbTransactionAuthMain)

	dbTransactionBlade, err := r.dbBlade.Begin()
	if err != nil {
		logger.DPanic("failed to begin Blade DB transaction", zap.Error(err),
			zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(dbTransactionBlade *sql.Tx) {
		logger.Debug("the Blade DB transaction was rollback", zap.String(requestIDKey, requestID))
		err := dbTransactionBlade.Rollback()
		if err != nil && err.Error() != messageTransactionHasAlreadyBeenCommittedOrRolledBack {
			logger.DPanic("failed to rollback Blade DB transaction", zap.Error(err),
				zap.String(requestIDKey, requestID))
		}
	}(dbTransactionBlade)

	// Lock tables

	_, err = dbTransactionAuthMain.Exec("/* postgreSQL query */\n" +
		"LOCK TABLE \"user\" IN SHARE ROW EXCLUSIVE MODE\n")
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	_, err = dbTransactionBlade.Exec("/* postgreSQL query */\n" +
		"LOCK TABLE confirmation_reset_password IN SHARE ROW EXCLUSIVE MODE\n")
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	// Check if an user exists with this email address

	loadUserID, err := dbTransactionAuthMain.Query("/* postgreSQL query */\n"+
		"SELECT\n"+
		"    \"user\".\"id\"\n"+
		"FROM\n"+
		"    \"user\"\n"+
		"WHERE\n"+
		"    \"user\".email = $1\n"+
		"    AND \"user\".deleted_at IS NULL\n"+
		"LIMIT 1\n",
		user.Email)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(loadUserID *sql.Rows) {
		if err := loadUserID.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(loadUserID)

	for loadUserID.Next() {
		err = loadUserID.Scan(&userID)
		if err != nil {
			logger.DPanic("failed to scan the data", zap.Error(err), zap.String(requestIDKey, requestID))
			return model.User{}, err
		}
	}

	if userID == 0 {
		return model.User{}, authorization.ErrorUserNotFound
	}

	// Updating the confirmation key status to "Deleted"

	deleteConfirmationKey, err := dbTransactionBlade.Prepare("/* postgreSQL query */\n" +
		"UPDATE\n" +
		"    confirmation_reset_password\n" +
		"SET\n" +
		"    deleted_at = NOW()\n" +
		"WHERE\n" +
		"    confirmation_reset_password.\"id\" = $1\n")
	if err != nil {
		logger.DPanic("failed to prepare the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(deleteConfirmationKey *sql.Stmt) {
		if err := deleteConfirmationKey.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(deleteConfirmationKey)

	_, err = deleteConfirmationKey.Exec(confirmationKeyID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	// Generate and update user password in the Auth DB

	user.Password = generator.StringRand(16, 16, false)
	hs := sha256.New()
	_, err = hs.Write([]byte(user.Password + r.config.CryptoSalt))
	if err != nil {
		logger.DPanic("failed to generate a new password", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	hashedPassword := hex.EncodeToString(hs.Sum(nil))

	updateUserPassword, err := dbTransactionAuthMain.Prepare("/* postgreSQL query */\n" +
		"UPDATE\n" +
		"    \"user\"\n" +
		"SET\n" +
		"    updated_at = NOW(),\n" +
		"    \"password\" = $1\n" +
		"WHERE\n" +
		"    \"user\".\"id\" = $2\n")
	if err != nil {
		logger.DPanic("failed to prepare the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}
	defer func(updateUserPassword *sql.Stmt) {
		if err := updateUserPassword.Close(); err != nil {
			logger.DPanic("failed to close result set", zap.Error(err), zap.String(requestIDKey, requestID))
		}
	}(updateUserPassword)

	_, err = updateUserPassword.Exec(hashedPassword, userID)
	if err != nil {
		logger.DPanic("failed to exec the query", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	// Transactions Commit

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		logger.DPanic("failed to commit to the AuthMain DB", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		logger.DPanic("failed to commit to the Blade DB", zap.Error(err), zap.String(requestIDKey, requestID))
		return model.User{}, err
	}

	return user, nil
}
