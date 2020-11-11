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
		PublicSessionID string `json:"id"`
		IssuedAt        int64  `json:"iat"`
	}
}

type headers struct {
	SigningAlgorithm string `json:"alg"`
	Type             string `json:"typ"`
}

type payload struct {
	Issuer          string `json:"iss"`
	Subject         string `json:"sub"`
	Purpose         string `json:"purpose"`
	PublicSessionID string `json:"id"`
	IssuedAt        int64  `json:"iat"`
}
