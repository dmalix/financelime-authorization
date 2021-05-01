/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"html"
	"net"
	"regexp"
	"strconv"
)

type inviteCode struct {
	ID          int64
	UserID      int64
	NumberLimit int
	Value       string
}

/*
	Create a new user
		----------------
		Return:
			confirmationID int64
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
				PROPS_EMAIL:              param user.Email is not valid
				PROPS_INVITE:             parap user.InviteCode is not valid
				PROPS_LANG:               param user.Language is not valid
				PROPS_REMOTE_ADDR:        param remoteAddr is not valid
				PROPS_LINK_KEY:           param linkKey is not valid
				USER_ALREADY_EXIST:       a user with the email you specified already exists
				INVITE_NOT_EXIST_EXPIRED: the invite code does not exist or is expired
				INVITE_LIMIT:             the limit for issuing this invite code has been exhausted
*/
// Related interfaces:
//	packages/authorization/domain/user.go
func (repo *Repo) CreateUser(user *models.User,
	remoteAddr, linkKey string, inviteCodeRequired bool) (int64, error) {

	type incomingProps struct {
		email              string
		inviteCode         string
		language           string
		remoteAddr         string
		inviteCodeRequired bool
		linkKey            string
	}

	var (
		props                   incomingProps
		dbTransactionAuthMain   *sql.Tx
		dbTransactionBlade      *sql.Tx
		dbRowsAuthMaster        *sql.Rows
		dbRowsBlade             *sql.Rows
		paramValueRegexp        *regexp.Regexp
		userID                  int64
		confirmationID          int64
		inviteCode              inviteCode
		inviteCodeReservedID    int64
		countInviteCodeIssued   int
		countInviteCodeReserved int
		inviteCodesIsRunOut     bool
		remoteAddrSource        net.IP
		remoteAddrResult        string
		err                     error
		errLabel                string
	)

	// Check props

	props.inviteCodeRequired = inviteCodeRequired

	props.email = html.EscapeString(user.Email)
	if len(props.email) <= 2 || len(props.email) > 255 {
		return confirmationID,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_EMAIL", "param user.Email is not valid", props.email))
	}

	props.inviteCode = html.EscapeString(user.InviteCode)
	paramValueRegexp = regexp.MustCompile(`^[0-9a-zA-Z_-]{3,16}$`)
	if !paramValueRegexp.MatchString(props.inviteCode) {
		if props.inviteCodeRequired {
			return confirmationID,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					"PROPS_INVITE", "param user.InviteCode is not valid", props.inviteCode))
		}
	}

	props.language = html.EscapeString(user.Language)
	paramValueRegexp = regexp.MustCompile(`^[ru|en]{2}$`)
	if !paramValueRegexp.MatchString(props.language) {
		return confirmationID,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_LANG", "param user.Language is not valid", props.language))
	}

	props.remoteAddr = html.EscapeString(remoteAddr)
	remoteAddrSource = net.ParseIP(props.remoteAddr)
	if remoteAddrSource == nil {
		return confirmationID,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_REMOTE_ADDR", "param remoteAddr is not valid", props.remoteAddr))
	}
	props.remoteAddr = remoteAddrSource.String()

	props.linkKey = html.EscapeString(linkKey)
	paramValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(props.linkKey) {
		return confirmationID,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_LINK_KEY", "param linkKey is not valid", props.linkKey))
	}

	// Begin the transaction

	dbTransactionAuthMain, err = repo.dbAuthMain.Begin()
	if err != nil {
		errLabel = "W0wfephh"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionAuthMain.Rollback()

	dbTransactionBlade, err = repo.dbBlade.Begin()
	if err != nil {
		errLabel = "FSvBG7Dr"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}
	//noinspection GoUnhandledErrorResult
	defer dbTransactionBlade.Rollback()

	// Lock tables

	_, err = dbTransactionAuthMain.Exec(`
		LOCK TABLE "user",
		invite_code,
		invite_code_issued IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		errLabel = "AA21lFGV"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	_, err = dbTransactionBlade.Exec(`
		LOCK TABLE confirmation_create_new_user,
		invite_code_reserved IN SHARE ROW EXCLUSIVE MODE`)
	if err != nil {
		errLabel = "KThpwB0c"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	if len(props.inviteCode) > 0 {

		// Check if this Invite Code exists

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
			return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		for dbRowsAuthMaster.Next() {
			err = dbRowsAuthMaster.Scan(&inviteCode.ID, &inviteCode.NumberLimit, &inviteCode.UserID)
			if err != nil {
				errLabel = "cWqgt3VB"
				return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
			}
		}

		if inviteCode.ID == 0 { // The invite code does not exist or is expired
			return confirmationID,
				errors.New(fmt.Sprintf("%s:%s[%s]",
					"INVITE_NOT_EXIST_EXPIRED", "the invite code does not exist or is expired", props.inviteCode))
		}

		// Check the limit for this invite code, including the reservation

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
				inviteCode.ID)
		if err != nil {
			errLabel = "P4BJAxNp"
			return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		for dbRowsAuthMaster.Next() {
			err = dbRowsAuthMaster.Scan(&countInviteCodeIssued)
			if err != nil {
				errLabel = "qooV4YZa"
				return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
				inviteCode.ID)
		if err != nil {
			errLabel = "K8bddqeW"
			return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		for dbRowsBlade.Next() {
			err = dbRowsBlade.Scan(&countInviteCodeReserved)
			errLabel = "exm38bTK"
			if err != nil {
				return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
			}
		}

		if (countInviteCodeIssued + countInviteCodeReserved) >= inviteCode.NumberLimit {
			inviteCodesIsRunOut = true

			if props.inviteCodeRequired { // the limit for issuing this invite code has been exhausted
				return confirmationID,
					errors.New(
						fmt.Sprintf(
							"%s:%s[inviteCode.NumberLimit=%s, countInviteCodeIssued=%s, countInviteCodeReserved=%s]",
							"INVITE_LIMIT", "the limit for issuing this invite code has been exhausted",
							strconv.Itoa(inviteCode.NumberLimit),
							strconv.Itoa(countInviteCodeIssued),
							strconv.Itoa(countInviteCodeReserved)))
			}
		}
	}

	// Check if an user exists with this email, including new users pending confirmation

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
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRowsAuthMaster.Next() {
		err = dbRowsAuthMaster.Scan(&userID)
		errLabel = "ygw0wRNX"
		if err != nil {
			return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRowsBlade.Next() {
		err = dbRowsBlade.Scan(&confirmationID)
		if err != nil {
			errLabel = "f8GLmoWc"
			return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	if userID != 0 || confirmationID != 0 { // a user with the email you specified already exists
		return confirmationID,
			errors.New(
				fmt.Sprintf(
					"%s:%s[userID=%d, confirmationID=%d]",
					"USER_ALREADY_EXIST", "a user with the email you specified already exists",
					userID,
					confirmationID))
	}

	// Add a new record for the user pending confirmation

	err =
		dbTransactionBlade.QueryRow(`
           	INSERT INTO confirmation_create_new_user 
           				( created_at, email, "language", link_key, remote_addr, expires_at )
			VALUES
				( NOW( ), $1, $2, $3, $4, NOW( ) + INTERVAL '1440 minute' ) RETURNING "id"`,
			props.email, props.language, props.linkKey, remoteAddrResult).Scan(&confirmationID)
	if err != nil {
		errLabel = "tC7ftRAS"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	if len(props.inviteCode) > 0 && !inviteCodesIsRunOut {

		// Add a new record for reserve the invite code for the user

		err =
			dbTransactionBlade.QueryRow(`
           		INSERT INTO invite_code_reserved ( created_at, invite_code_id, email, confirmation_id )
					VALUES
					( NOW( ), $1, $2, $3 ) RETURNING "id"`,
				inviteCode.ID, props.email, confirmationID).Scan(&inviteCodeReservedID)
		if err != nil {
			errLabel = "MANT4no8"
			return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	// Transactions Commit

	err = dbTransactionAuthMain.Commit()
	if err != nil {
		errLabel = "dnG1foyV"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	err = dbTransactionBlade.Commit()
	if err != nil {
		errLabel = "Dv3qdcSW"
		return confirmationID, errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	return confirmationID, nil
}
