package repository

import (
	"errors"
	"fmt"
	"html"
)

/*
	Delete the session
		----------------
		Return:
			err error  - system error code
*/
func (r *Repository) DeleteSession(userID int64, publicSessionID string) error {

	var (
		err      error
		errLabel string
	)

	_, err =
		r.dbAuthMain.Exec(`
	   		UPDATE 	"session" 
			SET 	deleted_at = NOW( )
			WHERE 	"session".user_id = $1 AND
					"session".public_id = $2`,
			userID, html.EscapeString(publicSessionID))
	if err != nil {
		errLabel = "3R19ASHZ"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	return nil
}
