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
	"github.com/dmalix/financelime-authorization/utils/trace"
	"hash"
)

/*
	Get User by Refresh Token
		----------------
		Return:
			user models.User
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
			    USER_NOT_FOUND: User is not found
*/

func (r *Repository) GetUserByRefreshToken(refreshToken string) (models.User, error) {

	var (
		user               models.User
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
		err = dbRows.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
		}
	}

	if user.ID == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"USER_NOT_FOUND", "User is not found", hashedRefreshToken))
	}

	return user, nil
}
