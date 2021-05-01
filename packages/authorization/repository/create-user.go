/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"net"
	"regexp"
	"strconv"
)

type inviteCodeRecord struct {
	ID          int64
	UserID      int64
	NumberLimit int
	Value       string
}

/*
	Create a new user
		----------------
		Return:
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
				PROPS_EMAIL:              the email param is not valid
				PROPS_LANG:               the propsUser.Language param is not valid
				PROPS_INVITE:             the propsInviteCode param is not valid
				PROPS_REMOTE_ADDR:        the propsRemoteAddr param is not valid
				PROPS_CONFIRMATION_KEY:   the propsConfirmationKey param is not valid
				USER_ALREADY_EXIST:       a user with the email you specified already exists
				INVITE_NOT_EXIST_EXPIRED: the invite code does not exist or is expired
				INVITE_LIMIT:             the limit for issuing this invite code has been exhausted
*/
// Related interfaces:
//	packages/authorization/domain.go
func (r *Repository) CreateUser(email, language, inviteCode, remoteAddr, confirmationKey string, inviteCodeRequired bool) error {

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
		inviteCodeRecord        inviteCodeRecord
		inviteCodeReservedID    int64
		countInviteCodeIssued   int
		countInviteCodeReserved int
		inviteCodesIsRunOut     bool
		confirmationID          int64
		remoteAddrSource        net.IP
		remoteAddrResult        string
		err                     error
		errLabel                string
	)

	// Check props
	// -----------

	props.inviteCodeRequired = inviteCodeRequired

	props.email = html.EscapeString(email)
	if len(props.email) <= 2 || len(props.email) > 255 {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_EMAIL", "param email is not valid", props.email))
	}

	props.inviteCode = html.EscapeString(inviteCode)
	paramValueRegexp = regexp.MustCompile(`^[0-9a-zA-Z_-]{3,16}$`)
	if !paramValueRegexp.MatchString(props.inviteCode) {
		if props.inviteCodeRequired {
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_INVITE", "param inviteCode is not valid", props.inviteCode))
		}
	}

	props.language = html.EscapeString(language)
	paramValueRegexp = regexp.MustCompile(`^[ru|en]{2}$`)
	if !paramValueRegexp.MatchString(props.language) {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_LANG", "param language is not valid", props.language))
	}

	props.remoteAddr = html.EscapeString(remoteAddr)
	remoteAddrSource = net.ParseIP(props.remoteAddr)
	if remoteAddrSource == nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_REMOTE_ADDR", "param remoteAddr is not valid", props.remoteAddr))
	}
	props.remoteAddr = remoteAddrSource.String()

	props.confirmationKey = html.EscapeString(confirmationKey)
	paramValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(props.confirmationKey) {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_CONFIRMATION_KEY", "param confirmationKey is not valid", props.confirmationKey))
	}

	// Begin the transaction
	// ---------------------

	dbTransactionAuthMain, err = r.dbAuthMain.Begin()
	if err != nil {
		errLabel = "W0wfephh"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionAuthMain.Rollback()

	dbTransactionBlade, err = r.dbBlade.Begin()
	if err != nil {
		errLabel = "FSvBG7Dr"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
		errLabel = "AA21lFGV"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	_, err = dbTransactionBlade.Exec(`
		LOCK TABLE confirmation_create_new_user,
		invite_code_reserved IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		errLabel = "KThpwB0c"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
			errLabel = "Chl5xLDp"
			return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		for dbRowsAuthMaster.Next() {
			err = dbRowsAuthMaster.Scan(&inviteCodeRecord.ID, &inviteCodeRecord.NumberLimit, &inviteCodeRecord.UserID)
			if err != nil {
				errLabel = "cWqgt3VB"
				return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
			}
		}

		if inviteCodeRecord.ID == 0 { // The invite code does not exist or is expired
			return errors.New(fmt.Sprintf("%s:%s[%s]",
				"INVITE_NOT_EXIST_EXPIRED", "the invite code does not exist or is expired", props.inviteCode))
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
				inviteCodeRecord.ID)
		if err != nil {
			errLabel = "P4BJAxNp"
			return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		for dbRowsAuthMaster.Next() {
			err = dbRowsAuthMaster.Scan(&countInviteCodeIssued)
			if err != nil {
				errLabel = "qooV4YZa"
				return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
				inviteCodeRecord.ID)
		if err != nil {
			errLabel = "K8bddqeW"
			return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		for dbRowsBlade.Next() {
			err = dbRowsBlade.Scan(&countInviteCodeReserved)
			errLabel = "exm38bTK"
			if err != nil {
				return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
			}
		}

		if (countInviteCodeIssued + countInviteCodeReserved) >= inviteCodeRecord.NumberLimit {
			inviteCodesIsRunOut = true

			if props.inviteCodeRequired { // the limit for issuing this invite code has been exhausted
				return errors.New(
					fmt.Sprintf(
						"%s:%s[inviteCode.NumberLimit=%s, countInviteCodeIssued=%s, countInviteCodeReserved=%s]",
						"INVITE_LIMIT", "the limit for issuing this invite code has been exhausted",
						strconv.Itoa(inviteCodeRecord.NumberLimit),
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
		errLabel = "sKc1YXnv"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRowsAuthMaster.Next() {
		err = dbRowsAuthMaster.Scan(&userID)
		errLabel = "ygw0wRNX"
		if err != nil {
			return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
		errLabel = "JJkxUbO7"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRowsBlade.Next() {
		err = dbRowsBlade.Scan(&confirmationID)
		if err != nil {
			errLabel = "f8GLmoWc"
			return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	if userID != 0 || confirmationID != 0 { // a user with the email you specified already exists
		return errors.New(
			fmt.Sprintf(
				"%s:%s[userID=%d, confirmationID=%d]",
				"USER_ALREADY_EXIST", "a user with the email you specified already exists",
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
		errLabel = "tC7ftRAS"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	if len(props.inviteCode) > 0 && !inviteCodesIsRunOut {

		// Add a new record for reserve the invite code for the user
		// ---------------------------------------------------------

		err =
			dbTransactionBlade.QueryRow(`
           		INSERT INTO invite_code_reserved ( created_at, invite_code_id, email, confirmation_id )
					VALUES
					( NOW( ), $1, $2, $3 ) RETURNING "id"`,
				inviteCodeRecord.ID, props.email, confirmationID).Scan(&inviteCodeReservedID)
		if err != nil {
			errLabel = "MANT4no8"
			return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	// Transactions Commit
	// -------------------

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		errLabel = "dnG1foyV"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		errLabel = "Dv3qdcSW"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	return nil
}
