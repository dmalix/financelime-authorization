/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"github.com/dmalix/financelime-rest-api/utils/random"
	"html"
	"regexp"
)

/*
	Confirm user email
		---------
		Return:
			user  models.User
			password string
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
			        ------------------------------------------------
			        PROPS_CONFIRMATION_KEY: The confirmationKey param is not valid
			        CONFIRMATION_KEY_NOT_FOUND_EXPIRED: The confirmation key hasn't found or expired.
			        CONFIRMATION_KEY_ALREADY_CONFIRMED: The user email is already confirmed.
*/
// Related interfaces:
//	packages/authorization/domain.go
func (r *Repository) ConfirmUserEmail(confirmationKey string) (models.User, error) {

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
		user                  models.User
	)

	propsValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !propsValueRegexp.MatchString(confirmationKey) {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_CONFIRMATION_KEY", "the confirmationKey param is not valid", html.EscapeString(confirmationKey)))
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
		err = dbRows.Scan(&confirmationID, &user.Email, &user.Language)
		if err != nil {
			errLabel = "h6PjQzPW"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}
	}

	if confirmationID == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"CONFIRMATION_KEY_NOT_FOUND_EXPIRED",
				"the confirmation key hasn't found or expired",
				html.EscapeString(confirmationKey)))
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
			user.Email)
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
			errors.New(fmt.Sprintf("%s:%s[%d]",
				"CONFIRMATION_KEY_ALREADY_CONFIRMED", "The user email is already confirmed.", userID))
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

	user.Password = random.StringRand(16, 16, false)
	hs := sha256.New()
	_, err = hs.Write([]byte(user.Password + r.cryptoSalt))
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
			user.Email, user.Language, hashedPassword).
			Scan(&user.ID)
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
				inviteCodeID, user.ID).
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
