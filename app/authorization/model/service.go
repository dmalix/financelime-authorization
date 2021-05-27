package model

type ServiceSignUpParam struct {
	Email      string
	Language   string
	InviteCode string
}

type ServiceCreateAccessTokenParam struct {
	Email     string
	Password  string
	ClientID  string
	UserAgent string
	Device    Device
}

type ServiceAccessTokenReturn struct {
	PublicSessionID string
	AccessJWT       string
	RefreshJWT      string
}

type ServiceRevokeRefreshTokenParam struct {
	EncryptedUserData []byte
	PublicSessionID   string
}
