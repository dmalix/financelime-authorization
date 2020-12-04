package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
)

/*
	Update the session
		----------------
		Return:
			error  - system or domain error code (format DOMAIN_ERROR_CODE:description[details]):
				------------------------------------------------
			    PROPS_REMOTE_ADDR: The remoteAddr param is not valid
				SESSION_NOT_FOUND: The case (the session + hashedRefreshToken) does not exist
*/

func (r *Repository) UpdateSession(publicSessionID, refreshToken, remoteAddr string) error {

	var (
		err                error
		errLabel           string
		remoteAddrSource   net.IP
		hashedRefreshToken string
		result             string
	)

	remoteAddrSource = net.ParseIP(remoteAddr)
	if remoteAddrSource == nil {
		return errors.New(fmt.Sprintf("%s:%s[%s]",
			"PROPS_REMOTE_ADDR", "The remoteAddr param is not valid", remoteAddr))
	}
	remoteAddr = remoteAddrSource.String()

	hs := sha256.New()
	_, err = hs.Write([]byte(
		refreshToken +
			r.config.CryptoSalt))
	if err != nil {
		errLabel = "Eeph6tho"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	hashedRefreshToken = hex.EncodeToString(hs.Sum(nil))

	err =
		r.dbAuthMain.QueryRow(`
           	UPDATE "session" 
			SET updated_at = NOW( ),
				remote_addr = $1,
				hashed_refresh_token = $2 
			WHERE
				"session".public_id = $3 AND  
				"session".hashed_refresh_token = $4 AND 
				"session".deleted_at IS NULL
			RETURNING "session".public_id`,
			remoteAddr, hashedRefreshToken, publicSessionID, hashedRefreshToken).Scan(&result)
	if err != nil {
		errLabel = "7AftRtCS"
		return errors.New(fmt.Sprintf("%s:[%s]", errLabel, err))
	}

	if len(result) == 0 {
		return errors.New(fmt.Sprintf("%s:%s[%s + %s]",
			"SESSION_NOT_FOUND", "the case (the session + hashedRefreshToken) does not exist",
			publicSessionID, hashedRefreshToken))
	}

	return nil
}
