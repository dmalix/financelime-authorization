/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"html"
	"net"
	"regexp"
)

/*
	Create a new user
		----------------
		Return:
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
				PROPS_EMAIL:              the email param is not valid
				PROPS_REMOTE_ADDR:        the propsRemoteAddr param is not valid
				PROPS_CONFIRMATION_KEY:   the propsConfirmationKey param is not valid
				USER_NOT_FOUND:           a user with the email specified not found
*/
// Related interfaces:
//	packages/authorization/domain.go
func (r *Repository) RequestUserPasswordReset(email, remoteAddr, confirmationKey string) (models.User, error) {

	type incomingProps struct {
		email           string
		remoteAddr      string
		confirmationKey string
	}

	var (
		props            incomingProps
		user             models.User
		query            string
		dbRows           *sql.Rows
		paramValueRegexp *regexp.Regexp
		remoteAddrSource net.IP
		confirmationID   int64
		err              error
		errLabel         string
	)

	// Check props
	// -----------

	props.email = html.EscapeString(email)
	if len(props.email) <= 2 || len(props.email) > 255 {
		return user, errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_EMAIL", "param email is not valid", props.email))
	}

	props.remoteAddr = html.EscapeString(remoteAddr)
	remoteAddrSource = net.ParseIP(props.remoteAddr)
	if remoteAddrSource == nil {
		return user, errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_REMOTE_ADDR", "param remoteAddr is not valid", props.remoteAddr))
	}
	props.remoteAddr = remoteAddrSource.String()

	props.confirmationKey = html.EscapeString(confirmationKey)
	paramValueRegexp = regexp.MustCompile(`^[abcefghijkmnopqrtuvwxyz23479]{16}$`)
	if !paramValueRegexp.MatchString(props.confirmationKey) {
		return user, errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_CONFIRMATION_KEY", "param confirmationKey is not valid", props.confirmationKey))
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

	dbRows, err = r.dbAuthRead.Query(query, email)
	if err != nil {
		errLabel = "FiCuhv2g"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	for dbRows.Next() {
		err = dbRows.Scan(&user.ID, &user.Email, &user.Language)
		if err != nil {
			errLabel = "RzZHe52q"
			return user,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))

		}
	}

	if user.ID == 0 {
		return user,
			errors.New(fmt.Sprintf("%s:%s[%s]",
				"USER_NOT_FOUND", "User is not found", email))
	}

	// Add the new record about reset password
	// ----------------------------------------

	err =
		r.dbBlade.QueryRow(`
           	INSERT 
           	INTO confirmation_reset_password (created_at, email, "language", confirmation_key, remote_addr, expires_at)
           		VALUES (NOW(), $1, $2, $3, $4, NOW() + interval '15 minute')
            RETURNING "id"`,
			email, user.Language, confirmationKey, remoteAddr).
			Scan(&confirmationID)
	if err != nil {
		errLabel = "d8A6m3WF"
		return user,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	return user, nil
}
