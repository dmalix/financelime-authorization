package repository

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-authorization/utils/trace"
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
		err error
	)

	_, err =
		r.dbAuthMain.Exec(`
	   		UPDATE 	"session" 
			SET 	deleted_at = NOW( )
			WHERE 	"session".user_id = $1 AND
					"session".public_id = $2`,
			userID, html.EscapeString(publicSessionID))
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", trace.GetCurrentPoint(), err))
	}

	return nil
}
