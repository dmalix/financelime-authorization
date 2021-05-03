/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/random"
	"github.com/dmalix/financelime-authorization/utils/trace"
	"hash"
	"html"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func NewPostgreDB(c ConfigPostgreDB) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s user=%s password=%s",
			c.Host, c.Port, c.SSLMode, c.DBName, c.User, c.Password))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s [%s]",
			trace.GetCurrentPoint(),
			"Failed to open DB connection",
			err.Error()))
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s [%s]",
			trace.GetCurrentPoint(),
			"Failed to ping to DB",
			err.Error()))
	}

	return db, nil
}

func NewRepository(
	config ConfigRepository,
	dbAuthMain,
	dbAuthRead,
	dbBlade *sql.DB) *repository {
	return &repository{
		config:     config,
		dbAuthMain: dbAuthMain,
		dbAuthRead: dbAuthRead,
		dbBlade:    dbBlade,
	}
}

func (r *repository) createUser(param repoCreateUserParam) error {

	type incomingProps struct {
		email              string
		inviteCode         string
		language           string
		remoteAddr         string
		inviteCodeRequired bool
		confirmationKey    string
	}

	var (
		props                   incomingProps
		dbTransactionAuthMain   *sql.Tx
		dbTransactionBlade      *sql.Tx
		dbRowsAuthMaster        *sql.Rows
		dbRowsBlade             *sql.Rows
		paramValueRegexp        *regexp.Regexp
		userID                  int64
		inviteCode              inviteCodeRecord
		inviteCodeReservedID    int64
		countInviteCodeIssued   int
		countInviteCodeReserved int
		inviteCodesIsRunOut     bool
		confirmationID          int64
		remoteAddrSource        net.IP
		remoteAddrResult        string
		err                     error
	)

	// Check props
	// -----------

	props.inviteCodeRequired = param.inviteCodeRequired

	props.email = html.EscapeString(param.email)
	if len(props.email) <= 2 || len(props.email) > 255 {
		return errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamEmail, props.email))
	}

	props.inviteCode = html.EscapeString(param.inviteCode)
	paramValueRegexp = regexp.MustCompile(`^[0-9a-zA-Z_-]{3,16}$`)
	if !paramValueRegexp.MatchString(props.inviteCode) {
		if props.inviteCodeRequired {
			return errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeBadParamInvite, props.inviteCode))
		}
	}

	props.language = html.EscapeString(param.language)
	paramValueRegexp = regexp.MustCompile(`^[ru|en]{2}$`)
	if !paramValueRegexp.MatchString(props.language) {
		return errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamLang, props.language))
	}

	props.remoteAddr = html.EscapeString(param.remoteAddr)
	remoteAddrSource = net.ParseIP(props.remoteAddr)
	if remoteAddrSource == nil {
		return errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamRemoteAddr, props.remoteAddr))
	}
	props.remoteAddr = remoteAddrSource.String()

	props.confirmationKey = html.EscapeString(string(param.confirmationKey))
	paramValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(props.confirmationKey) {
		return errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamConfirmationKey, props.confirmationKey))
	}

	// Begin the transaction
	// ---------------------

	dbTransactionAuthMain, err = r.dbAuthMain.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionAuthMain.Rollback()

	dbTransactionBlade, err = r.dbBlade.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionBlade.Rollback()

	// Lock tables
	// -----------

	_, err = dbTransactionAuthMain.Exec(`
		LOCK TABLE "user",
		invite_code,
		invite_code_issued IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	_, err = dbTransactionBlade.Exec(`
		LOCK TABLE confirmation_create_new_user,
		invite_code_reserved IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	if len(props.inviteCode) > 0 {

		// Check if this Invite Code exists
		// --------------------------------

		dbRowsAuthMaster, err =
			dbTransactionAuthMain.Query(`
			SELECT
				invite_code."id",
				invite_code.number_limit,
				"user"."id" AS user_id 
			FROM
				invite_code
				INNER JOIN "user" ON invite_code.user_id = "user"."id" 
			WHERE
				invite_code."value" = $1 
				AND "user".deleted_at IS NULL 
				AND invite_code.deleted_at IS NULL 
				AND invite_code.expires_at > NOW( ) 
				LIMIT 1`,
				props.inviteCode)
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}

		for dbRowsAuthMaster.Next() {
			err = dbRowsAuthMaster.Scan(&inviteCode.id, &inviteCode.numberLimit, &inviteCode.userID)
			if err != nil {
				return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
			}
		}

		if inviteCode.id == 0 { // The invite code does not exist or is expired
			return errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeInviteNotFound, props.inviteCode))
		}

		// Check the limit for this invite code, including the reservation
		// ---------------------------------------------------------------

		dbRowsAuthMaster, err =
			dbTransactionAuthMain.Query(`
			SELECT COUNT
				( invite_code_issued."id" ) 
			FROM
				invite_code
				INNER JOIN "user" ON invite_code.user_id = "user"."id"
				INNER JOIN invite_code_issued ON invite_code."id" = invite_code_issued.invite_code_id 
			WHERE
				invite_code."id" = $1 
				AND "user".deleted_at IS NULL 
				AND invite_code_issued.deleted_at IS NULL 
				AND invite_code.deleted_at IS NULL 
				AND invite_code.expires_at > NOW( )`,
				inviteCode.id)
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}

		for dbRowsAuthMaster.Next() {
			err = dbRowsAuthMaster.Scan(&countInviteCodeIssued)
			if err != nil {
				return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
			}
		}

		dbRowsBlade, err =
			dbTransactionBlade.Query(`
			SELECT COUNT
				( invite_code_reserved."id" ) 
			FROM
				invite_code_reserved
				INNER JOIN confirmation_create_new_user 
					ON invite_code_reserved.email = confirmation_create_new_user.email 
			WHERE
				invite_code_reserved.invite_code_id = $1 
				AND invite_code_reserved.deleted_at IS NULL 
				AND confirmation_create_new_user.deleted_at IS NULL 
				AND confirmation_create_new_user.expires_at > NOW( )`,
				inviteCode.id)
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}

		for dbRowsBlade.Next() {
			err = dbRowsBlade.Scan(&countInviteCodeReserved)
			if err != nil {
				return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
			}
		}

		if (countInviteCodeIssued + countInviteCodeReserved) >= inviteCode.numberLimit {
			inviteCodesIsRunOut = true

			if props.inviteCodeRequired { // the limit for issuing this invite code has been exhausted
				return errors.New(
					fmt.Sprintf(
						"%s:inviteCode.NumberLimit=%s,countInviteCodeIssued=%s,countInviteCodeReserved=%s",
						domainErrorCodeInviteHasEnded,
						strconv.Itoa(inviteCode.numberLimit),
						strconv.Itoa(countInviteCodeIssued),
						strconv.Itoa(countInviteCodeReserved)))
			}
		}
	}

	// Check if an user exists with this email, including new users pending confirmation
	// ---------------------------------------------------------------------------------

	dbRowsAuthMaster, err =
		dbTransactionAuthMain.Query(`
		SELECT 
			"user"."id" 
		FROM 
			"user" 
		WHERE 
			"user".email = $1 
			AND "user".deleted_at IS NULL 
			LIMIT 1`,
			props.email)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	for dbRowsAuthMaster.Next() {
		err = dbRowsAuthMaster.Scan(&userID)
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}
	}

	dbRowsBlade, err =
		dbTransactionBlade.Query(`
		SELECT
			confirmation_create_new_user."id" 
		FROM
			confirmation_create_new_user 
		WHERE
			confirmation_create_new_user.email = $1 
			AND confirmation_create_new_user.deleted_at IS NULL 
			AND confirmation_create_new_user.expires_at > NOW( ) 
		LIMIT 1`,
			props.email)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	for dbRowsBlade.Next() {
		err = dbRowsBlade.Scan(&confirmationID)
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}
	}

	if userID != 0 || confirmationID != 0 { // a user with the email you specified already exists
		return errors.New(
			fmt.Sprintf(
				"%s:userID=%d,confirmationID=%d",
				domainErrorCodeUserAlreadyExist,
				userID,
				confirmationID))
	}

	// Add a new record for the user pending confirmation
	// --------------------------------------------------

	err =
		dbTransactionBlade.QueryRow(`
           	INSERT INTO confirmation_create_new_user 
           				( created_at, email, "language", confirmation_key, remote_addr, expires_at )
			VALUES
				( NOW( ), $1, $2, $3, $4, NOW( ) + INTERVAL '1440 minute' ) RETURNING "id"`,
			props.email, props.language, props.confirmationKey, remoteAddrResult).Scan(&confirmationID)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	if len(props.inviteCode) > 0 && !inviteCodesIsRunOut {

		// Add a new record for reserve the invite code for the user
		// ---------------------------------------------------------

		err =
			dbTransactionBlade.QueryRow(`
           		INSERT INTO invite_code_reserved ( created_at, invite_code_id, email, confirmation_id )
					VALUES
					( NOW( ), $1, $2, $3 ) RETURNING "id"`,
				inviteCode.id, props.email, confirmationID).Scan(&inviteCodeReservedID)
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}
	}

	// Transactions Commit
	// -------------------

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	return nil
}

/*
	Confirm user email
		---------
		Return:
			user  models.User
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
			        ------------------------------------------------
			        PROPS_CONFIRMATION_KEY: The confirmationKey param is not valid
			        CONFIRMATION_KEY_NOT_FOUND_EXPIRED: The confirmation key hasn't found or expired.
			        CONFIRMATION_KEY_ALREADY_CONFIRMED: The user email is already confirmed.
*/
// Related interfaces:
//	packages/authorization/domain.go
func (r *repository) confirmUserEmail(confirmationKey string) (user, error) {

	var (
		dbTransactionAuthMain *sql.Tx
		dbTransactionBlade    *sql.Tx
		dbRows                *sql.Rows
		propsValueRegexp      *regexp.Regexp
		err                   error
		errLabel              string
		dbStmt                *sql.Stmt
		confirmationID        int64
		inviteCodeIssuedID    int64
		inviteCodeID          int64
		userID                int64
		user                  user
	)

	propsValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !propsValueRegexp.MatchString(string(confirmationKey)) {
		return user,
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeBadParamConfirmationKey, html.EscapeString(string(confirmationKey))))
	}

	// Check the confirmationKey in Database
	// -------------------------------------

	dbRows, err =
		r.dbBlade.Query(`
			SELECT
				confirmation_create_new_user."id",
				confirmation_create_new_user.email,
				confirmation_create_new_user."language"
			FROM
				confirmation_create_new_user
			WHERE
				confirmation_create_new_user.confirmation_key = $1 
				AND confirmation_create_new_user.deleted_at IS NULL 
				AND confirmation_create_new_user.expires_at > NOW( ) 
			ORDER BY
				confirmation_create_new_user."id" DESC 
				LIMIT 1`,
			confirmationKey)
	if err != nil {
		errLabel = "Dtv3CkDF"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	//noinspection GoUnhandledErrorResult
	defer dbRows.Close()

	for dbRows.Next() {
		err = dbRows.Scan(&confirmationID, &user.email, &user.language)
		if err != nil {
			errLabel = "h6PjQzPW"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	if confirmationID == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeConfirmationKeyNotFound,
				html.EscapeString(string(confirmationKey))))
	}

	// Begin the transaction
	// ---------------------

	dbTransactionAuthMain, err = r.dbAuthMain.Begin()
	if err != nil {
		errLabel = "oBfY4KaU"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionAuthMain.Rollback()

	dbTransactionBlade, err = r.dbBlade.Begin()
	if err != nil {
		errLabel = "UDg6e1a9"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionBlade.Rollback()

	// Lock tables
	// -----------

	_, err = dbTransactionAuthMain.Exec(`
		LOCK TABLE "user",
		invite_code,
		invite_code_issued IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		errLabel = "Tfdo686a"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	_, err = dbTransactionBlade.Exec(`
		LOCK TABLE confirmation_create_new_user,
		invite_code_reserved IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		errLabel = "Iru05LVw"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	// Check if an user exists with this email address
	// -----------------------------------------------

	dbRows, err =
		dbTransactionAuthMain.Query(`
			SELECT 
				"user"."id" 
			FROM 
				"user" 
			WHERE 
				"user".email = $1 
				AND "user".deleted_at IS NULL 
				LIMIT 1`,
			user.email)
	if err != nil {
		errLabel = "Dtv3CkDF"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	//noinspection GoUnhandledErrorResult
	defer dbRows.Close()

	for dbRows.Next() {
		err = dbRows.Scan(&userID)
		if err != nil {
			errLabel = "B95rVIkG"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	if userID != 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%d",
				domainErrorCodeConfirmationKeyAlreadyConfirmed, userID))
	}

	// Updating the confirmation key status to "Deleted"
	// -------------------------------------------------

	dbStmt, err =
		dbTransactionBlade.Prepare(`
	   		UPDATE confirmation_create_new_user 
			SET deleted_at = NOW( ) 
			WHERE
				confirmation_create_new_user."id" = $1`)
	if err != nil {
		errLabel = "fWTLD7vM"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	_, err = dbStmt.Exec(confirmationID)
	if err != nil {
		errLabel = "LgTumH8j"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	// Check if there was a reserve for an invite code
	// -----------------------------------------------

	dbRows, err =
		dbTransactionBlade.Query(`
	   		SELECT
				invite_code_reserved.invite_code_id 
			FROM
				invite_code_reserved 
			WHERE
				invite_code_reserved.confirmation_id = $1 
				AND invite_code_reserved.deleted_at IS NULL
				LIMIT 1`,
			confirmationID)
	if err != nil {
		errLabel = "oPT8P0oY"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRows.Next() {
		err = dbRows.Scan(&inviteCodeID)
		if err != nil {
			errLabel = "KdUlzN7C"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	if inviteCodeID > 0 {

		// Updating the invite code reserve to "Deleted"
		// ---------------------------------------------

		dbStmt, err =
			dbTransactionBlade.Prepare(`
	   		UPDATE invite_code_reserved 
			SET deleted_at = NOW( ) 
			WHERE
				invite_code_reserved.confirmation_id = $1`)
		if err != nil {
			errLabel = "RcppI0gB"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		_, err = dbStmt.Exec(confirmationID)
		if err != nil {
			errLabel = "rKAn8tLM"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	// Create the new confirmed user in the Auth DB
	// --------------------------------------------

	user.password = random.StringRand(16, 16, false)
	hs := sha256.New()
	_, err = hs.Write([]byte(user.password + r.config.CryptoSalt))
	if err != nil {
		errLabel = "Jhc7OYxi"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}
	hashedPassword := hex.EncodeToString(hs.Sum(nil))

	err =
		dbTransactionAuthMain.QueryRow(`
	   		INSERT INTO "user" ( created_at, email, "language", "password" )
			VALUES
				( NOW( ), $1, $2, $3 ) RETURNING "id"`,
			user.email, user.language, hashedPassword).
			Scan(&user.id)
	if err != nil {
		errLabel = "Annm39Zs"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	if inviteCodeID > 0 {

		// Linking the invite code to the user
		// -----------------------------------

		err =
			dbTransactionAuthMain.QueryRow(`
	   		INSERT INTO invite_code_issued ( created_at, invite_code_id, user_id )
			VALUES
				( NOW( ), $1, $2 ) RETURNING "id"`,
				inviteCodeID, user.id).
				Scan(&inviteCodeIssuedID)
		if err != nil {
			errLabel = "p3He5f9l"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	// Transactions Commit
	// -------------------

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		errLabel = "HI1ZZ1CP"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		errLabel = "BG1zZftV"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	return user, nil
}

func (r *repository) getUserByAuth(param repoGetUserByAuthParam) (user, error) {

	var (
		user           user
		dbRows         *sql.Rows
		query          string
		hs             hash.Hash
		hashedPassword string
		err            error
	)

	// Check Props
	// -----------

	if len(param.email) <= 2 || len(param.email) > 255 {
		return user,
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeBadParamEmail, param.email))
	}

	if len(param.password) <= 2 || len(param.password) > 255 {
		return user,
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeBadParamPassword, param.password))
	}

	// Get User
	// --------

	hs = sha256.New()
	_, err = hs.Write([]byte(param.password + r.config.CryptoSalt))
	if err != nil {
		return user, errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}
	hashedPassword = hex.EncodeToString(hs.Sum(nil))

	query = `
		SELECT
			"user"."id",
			"user".email,
			"user"."language"
		FROM
			"user"
		WHERE
			"user".email = $1 AND
			"user".password = $2 AND
			"user".deleted_at IS NULL  
			LIMIT 1`

	dbRows, err = r.dbAuthRead.Query(query, param.email, hashedPassword)
	if err != nil {
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	for dbRows.Next() {
		err = dbRows.Scan(&user.id, &user.email, &user.language)
		if err != nil {
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))

		}
	}

	if user.id == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeUserNotFound, param.password))
	}

	return user, nil
}

/*
	Get User by Refresh Token
		----------------
		Return:
			user models.User
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
			    USER_NOT_FOUND: User is not found
*/

func (r *repository) getUserByRefreshToken(refreshToken string) (user, error) {

	var (
		user               user
		dbRows             *sql.Rows
		query              string
		hs                 hash.Hash
		hashedRefreshToken string
		err                error
	)

	hs = sha256.New()
	_, err = hs.Write([]byte(
		refreshToken +
			r.config.CryptoSalt))
	if err != nil {
		return user, errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	hashedRefreshToken = hex.EncodeToString(hs.Sum(nil))

	query = `
		SELECT 	"user"."id", 
				"user".email, 
				"user"."language"
		FROM 	"session"
		INNER JOIN "user" ON "session".user_id = "user"."id"  
		WHERE	"session".hashed_refresh_token = $1 AND
				"session".deleted_at IS NULL AND
				"user".deleted_at IS NULL
		LIMIT 1`

	dbRows, err = r.dbAuthRead.Query(query, hashedRefreshToken)
	if err != nil {
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	for dbRows.Next() {
		err = dbRows.Scan(&user.id, &user.email, &user.language)
		if err != nil {
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}
	}

	if user.id == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"USER_NOT_FOUND", "User is not found", hashedRefreshToken))
	}

	return user, nil
}

func (r *repository) saveSession(param repoSaveSessionParam) error {

	var (
		err              error
		sessionID        int64
		deviceID         int64
		remoteAddrSource net.IP
	)

	remoteAddrSource = net.ParseIP(param.remoteAddr)
	if remoteAddrSource == nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_REMOTE_ADDR", "The remoteAddr param is not valid", param.remoteAddr))
	}
	param.remoteAddr = remoteAddrSource.String()

	hs := sha256.New()
	_, err = hs.Write([]byte(
		param.refreshToken +
			r.config.CryptoSalt))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			trace.GetCurrentPoint(),
			"Failed to generate hash for the Refresh Token",
			err))
	}

	hashedRefreshToken := hex.EncodeToString(hs.Sum(nil))

	err =
		r.dbAuthMain.QueryRow(`
	   		INSERT INTO "session" ( created_at, user_id, client_id, remote_addr, public_id, hashed_refresh_token )
			VALUES
				( NOW( ), $1, $2, $3, $4, $5 ) RETURNING "id"`,
			param.userID, param.clientID, param.remoteAddr, param.publicSessionID, hashedRefreshToken).
			Scan(&sessionID)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	err =
		r.dbAuthMain.QueryRow(`
	   		INSERT INTO device ( created_at, session_id, platform, height, width, "language", timezone, user_agent )
			VALUES
				( NOW( ), $1, $2, $3, $4, $5, $6, $7  ) RETURNING "id"`,
			sessionID,
			html.EscapeString(param.device.Platform),
			param.device.Height,
			param.device.Width,
			html.EscapeString(param.device.Language),
			html.EscapeString(param.device.Timezone),
			html.EscapeString(param.userAgent)).
			Scan(&deviceID)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	return nil
}

/*
	Update the session
		----------------
		Return:
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
			    PROPS_REMOTE_ADDR: The remoteAddr param is not valid
				SESSION_NOT_FOUND: The case (the session + hashedRefreshToken) does not exist
*/

func (r *repository) updateSession(param repoUpdateSessionParam) error {

	var (
		err    error
		result string
	)

	remoteAddrSource := net.ParseIP(param.remoteAddr)
	if remoteAddrSource == nil {
		return errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamRemoteAddr, param.remoteAddr))
	}
	remoteAddr := remoteAddrSource.String()

	hs := sha256.New()
	_, err = hs.Write([]byte(
		param.refreshToken +
			r.config.CryptoSalt))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	hashedRefreshToken := hex.EncodeToString(hs.Sum(nil))

	err =
		r.dbAuthMain.QueryRow(`
           	UPDATE "session" 
			SET updated_at = NOW( ),
				remote_addr = $1,
				hashed_refresh_token = $2 
			WHERE
				"session".public_id = $3 AND  
				"session".hashed_refresh_token = $4 AND 
				"session".deleted_at IS NULL
			RETURNING "session".public_id`,
			remoteAddr, hashedRefreshToken, param.publicSessionID, hashedRefreshToken).Scan(&result)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	if len(result) == 0 {
		return errors.New(fmt.Sprintf("%s:%s[%s + %s]",
			domainErrorSessionNotFound, "the case (the session + hashedRefreshToken) does not exist",
			param.publicSessionID, hashedRefreshToken))
	}

	return nil
}

func (r *repository) requestUserPasswordReset(param repoRequestUserPasswordResetParam) (user, error) {

	var (
		user             user
		query            string
		dbRows           *sql.Rows
		paramValueRegexp *regexp.Regexp
		remoteAddrSource net.IP
		confirmationID   int64
		err              error
	)

	// Check props
	// -----------

	email := html.EscapeString(param.email)
	if len(email) <= 2 || len(email) > 255 {
		return user, errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamEmail, email))
	}

	remoteAddr := html.EscapeString(param.remoteAddr)
	remoteAddrSource = net.ParseIP(remoteAddr)
	if remoteAddrSource == nil {
		return user, errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadParamRemoteAddr, remoteAddr))
	}
	remoteAddr = remoteAddrSource.String()

	confirmationKey := html.EscapeString(param.confirmationKey)
	paramValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(confirmationKey) {
		return user, errors.New(fmt.Sprintf("%s:%s",
			domainErrorCodeBadConfirmationKey, confirmationKey))
	}

	// Check if the user exists
	// ------------------------

	query = `
		SELECT
			"user"."id",
			"user".email,
			"user"."language"
		FROM
			"user"
		WHERE
			"user".email = $1 AND
			"user".deleted_at IS NULL  
			LIMIT 1`

	dbRows, err = r.dbAuthRead.Query(query, param.email)
	if err != nil {
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	for dbRows.Next() {
		err = dbRows.Scan(&user.id, &user.email, &user.language)
		if err != nil {
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))

		}
	}

	if user.id == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s",
				domainErrorCodeUserNotFound, param.email))
	}

	// Add the new record about reset password
	// ----------------------------------------

	err =
		r.dbBlade.QueryRow(`
           	INSERT 
           	INTO confirmation_reset_password (created_at, email, "language", confirmation_key, remote_addr, expires_at)
           		VALUES (NOW(), $1, $2, $3, $4, NOW() + interval '15 minute')
            RETURNING "id"`,
			param.email, user.language, param.confirmationKey, param.remoteAddr).
			Scan(&confirmationID)
	if err != nil {
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	return user, nil
}

func (r *repository) getListActiveSessions(userID int64) ([]session, error) {
	var (
		dbRows   *sql.Rows
		query    string
		sessions []session
		session  session
		err      error
	)

	query = strings.Replace(`
		SELECT
			"session".public_id,
		  	device.platform,	
			CASE 
				WHEN ( "session".updated_at IS NULL )
						THEN "session".created_at
						ELSE "session".updated_at
				END AS updated_at
		FROM
			"session"
		  	INNER JOIN device ON "session"."id" = device.session_id	
		WHERE
			(
				( "session".updated_at IS NULL AND "session".created_at > NOW( ) - INTERVAL 'LIFETIME SECOND' ) 
				OR ( "session".updated_at IS NOT NULL AND "session".updated_at > NOW( ) - INTERVAL 'LIFETIME SECOND' ) 
			) 
			AND "session".user_id = $1 
			AND "session".deleted_at IS NULL`,
		"LIFETIME", strconv.Itoa(r.config.JwtRefreshTokenLifetime), 2)

	dbRows, err = r.dbAuthRead.Query(query, userID)

	if err != nil {
		return sessions,
			errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	//noinspection GoUnhandledErrorResult
	defer dbRows.Close()

	for dbRows.Next() {
		err = dbRows.Scan(
			&session.PublicSessionID,
			&session.Platform,
			&session.UpdatedAt,
		)
		if err != nil {
			return sessions,
				errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *repository) deleteSession(param repoDeleteSessionParam) error {

	var (
		err error
	)

	_, err =
		r.dbAuthMain.Exec(`
	   		UPDATE 	"session" 
			SET 	deleted_at = NOW( )
			WHERE 	"session".user_id = $1 AND
					"session".public_id = $2`,
			param.userID, html.EscapeString(param.publicSessionID))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	return nil
}
