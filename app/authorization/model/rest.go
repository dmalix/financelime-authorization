package model

type SignUpRequest struct {
	// User email
	Email string `json:"email" validate:"required" example:"test.user@financelime.com"`
	// Invite code. Required depending on the setting of environment variable `AUTH_INVITE_CODE_REQUIRED`. The default is required.
	InviteCode string `json:"inviteCode" validate:"required" example:"testInviteCode"`
	// User language
	Language string `json:"language" validate:"required" example:"en"`
}

type CreateAccessTokenRequest struct {
	// User Email
	Email string `json:"email" validate:"required" example:"test.user@financelime.com"`
	// User Password
	Password string `json:"password" validate:"required" example:"qmhVXVC1%hVNa0Hcq"`
	// User Client ID
	ClientID string `json:"clientID" validate:"required" example:"PWA_v0.0.1"`

	Device Device `json:"device" validate:"required"`
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJmaW5hbmNlbGltZS5jb20iLCJzdWIiOiJBdXRob3JpemF0aW9uIiwicHVycG9zZSI6InJlZnJlc2giLCJzZXNzaW9uSUQiOiI2M2IyZjUyM2ZiZGQzMzFlZjQzM2U2NmU5NDZjYWQ2OTNkOTQ5MzdjMWUxNWNjMDI5YjFiNjE1YmExN2VjZWM3IiwidXNlckRhdGEiOiJpNHhTbDBmNXcrMGJSTE1KOURMVlp3NGJDZkpSSUFqQlhoN2VFZFlpNTV2L1QvVk1EK3RmNFNyK0NSV0ZnZEpoUkh2S1AyNnZGQ1AxZ05iOU4yejljMFRoYkRZNkFSdGt2WHkzMHJ3bTlDeXh0Vk1QdUUvRXh4UDdzaCs3MGVrbE5ObjdGS2xIIiwiaWF0IjoxNjIwNTIwNTg2fQ.4fd650daddded3a74c6fcfa28559d02c3ca6f32d55805b58ac88ccd302c5445e"`
}

type AccessTokenResponse struct {
	PublicSessionID string `json:"sessionID"`
	AccessJWT       string `json:"accessToken"`
	RefreshJWT      string `json:"refreshToken"`
}

type RevokeRefreshTokenRequest struct {
	PublicSessionID string `json:"sessionID" validate:"required" example:"f58f06a96b69083b7c4fb068faa6c8314af0636e44ecc710261abe1759b07755"`
}

type ResetUserPasswordRequest struct {
	Email string `json:"email" validate:"required" example:"test.user@financelime.com"`
}

/////////////////////////////////////////////////////////////

type CommonFailure struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"404 Not Found"`
}

type SignUpFailure400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" enums:"BAD_PARAMETERS" example:"BAD_PARAMETERS"`
}

type SignUpFailure409 struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" enums:"USER_ALREADY_EXIST,INVITE_NOT_FOUND,INVITE_HAS_ENDED" example:"USER_ALREADY_EXIST"`
}

type CreateAccessTokenFailure400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" enums:"BAD_PARAMETERS" example:"BAD_PARAMETERS"`
}

type CreateAccessTokenFailure404 struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" enums:"USER_NOT_FOUND" example:"USER_NOT_FOUND"`
}

type RefreshAccessTokenFailure400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" enums:"BAD_PARAMETERS, BAD_REFRESH_TOKEN" example:"BAD_PARAMETERS"`
}

type RefreshAccessTokenFailure404 struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" enums:"USER_NOT_FOUND" example:"USER_NOT_FOUND"`
}

type RevokeRefreshTokenFailure400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" enums:"BAD_PARAMETERS" example:"BAD_PARAMETERS"`
}

type RequestUserPasswordResetFailure400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" enums:"BAD_PARAMETERS, BAD_REFRESH_TOKEN" example:"BAD_PARAMETERS"`
}

type RequestUserPasswordResetFailure404 struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" enums:"USER_NOT_FOUND" example:"USER_NOT_FOUND"`
}
