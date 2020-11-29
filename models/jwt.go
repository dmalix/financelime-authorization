package models

type JwtData struct {
	Headers struct {
		SigningAlgorithm string `json:"alg"`
		Type             string `json:"typ"`
	}
	Payload struct {
		Issuer          string `json:"iss"`
		Subject         string `json:"sub"`
		Purpose         string `json:"purpose"`
		PublicSessionID string `json:"sessionID"`
		UserData        []byte `json:"userData"`
		IssuedAt        int64  `json:"iat"`
	}
}
