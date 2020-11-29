/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/models"
	"strconv"
	"strings"
)

/*
	   	Get a list of active sessions
	   		----------------
	   		Return:
				sessions []models.Session
	   			err error
*/
func (r *Repository) GetListActiveSessions(userID int64) ([]models.Session, error) {
	var (
		dbRows   *sql.Rows
		query    string
		sessions []models.Session
		session  models.Session
		err      error
		errLabel string
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
		errLabel = "lFAA21GV"
		return sessions,
			errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
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
			errLabel = "dqoZZ4fR"
			return sessions,
				errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
