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
	"github.com/dmalix/financelime-authorization/models"
	"hash"
	"log"
)

/*
	Get User by auth
		----------------
		Return:
			user models.User
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
			    PROPS_EMAIL: Email param is not valid.
			    PROPS_PASSWORD: Password param is not valid.
			    USER_NOT_FOUND: User is not found.
*/
func (r *Repository) GetUserByAuth(email, password string) (models.User, error) {

	var (
		user           models.User
		dbRows         *sql.Rows
		query          string
		hs             hash.Hash
		hashedPassword string
		err            error
		errLabel       string
	)

	// Check Props
	// -----------

	if len(email) <= 2 || len(email) > 255 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_EMAIL", "Email param is not valid", email))
	}

	if len(password) <= 2 || len(password) > 255 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"PROPS_PASSWORD", "Password param is not valid", password))
	}

	// Get User
	// --------

	hs = sha256.New()
	_, err = hs.Write([]byte(password + r.config.CryptoSalt))
	if err != nil {
		return user, errors.New(fmt.Sprintf("%s:[%s]", "XAacJE73", err))
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

	log.Println("=================================================================================================")
	log.Println("email")
	log.Println(email)
	log.Println("-------------------------------------------------------------------------------------------------")
	log.Println("hashedPassword")
	log.Println(hashedPassword)
	log.Println("=================================================================================================")

	dbRows, err = r.dbAuthRead.Query(query, email, hashedPassword)
	if err != nil {
		errLabel = "hv2FiCug"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRows.Next() {
		err = dbRows.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			errLabel = "ZHe5Rz2q"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))

		}
	}

	if user.ID == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"USER_NOT_FOUND", "User is not found", password))
	}

	return user, nil
}
