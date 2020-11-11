package repository

import (
	"errors"
	"fmt"
	"github.com/dmalix/financelime-rest-api/models"
	"html"
	"net"
)

/*
	Save the session
		----------------
		Return:
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
			    PROPS_REMOTE_ADDR: The remoteAddr param is not valid
*/
func (r *Repository) SaveSession(userID int64, publicSessionID, clientID, remoteAddr string, device models.Device) error {

	var (
		err              error
		sessionID        int64
		deviceID         int64
		remoteAddrSource net.IP
	)

	remoteAddrSource = net.ParseIP(remoteAddr)
	if remoteAddrSource == nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_REMOTE_ADDR", "The remoteAddr param is not valid", remoteAddr))
	}
	remoteAddr = remoteAddrSource.String()

	err =
		r.dbAuthMain.QueryRow(`
	   		INSERT INTO "session" ( created_at, user_id, client_id, remote_addr, public_id )
			VALUES
				( NOW( ), $1, $2, $3, $4  ) RETURNING "id"`,
			userID, clientID, remoteAddr, publicSessionID).
			Scan(&sessionID)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", "kM4BsYfY", err))
	}

	err =
		r.dbAuthMain.QueryRow(`
	   		INSERT INTO device ( created_at, session_id, platform, height, width, "language", timezone, user_agent )
			VALUES
				( NOW( ), $1, $2, $3, $4, $5, $6, $7  ) RETURNING "id"`,
			sessionID,
			html.EscapeString(device.Platform),
			device.Height,
			device.Width,
			html.EscapeString(device.Language),
			html.EscapeString(device.Timezone),
			html.EscapeString(device.UserAgent)).
			Scan(&deviceID)
	if err != nil {
		return errors.New(fmt.Sprintf("%s:[%s]", "wU3fcOtz", err))
	}

	return nil
}
