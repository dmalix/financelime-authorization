package jwt

type JWT interface {
	GenerateToken(publicSessionID string, privateBox []byte, tokenPurpose string, issuedAt ...int64) (string, error)
	VerifyToken(jwt string) (JsonWebToken, error)
}
