package jwt

type Jwt interface {
	GenerateToken(publicSessionID string, userData []byte, tokenPurpose string, issuedAt ...int64) (string, error)
	VerifyToken(jwt string) (JwtData, error)
}
