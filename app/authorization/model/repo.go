package model

type RepoSignUpParam struct {
	Email              string
	Language           string
	InviteCode         string
	ConfirmationKey    string
	InviteCodeRequired bool
}
type RepoSaveSessionParam struct {
	UserID          int64
	PublicSessionID string
	RefreshToken    string
	ClientID        string
	UserAgent       string
	Device          Device
}

type RepoGetUserByAuthParam struct {
	Email    string
	Password string
}

type RepoDeleteSessionParam struct {
	UserID          int64
	PublicSessionID string
}

type RepoUpdateSessionParam struct {
	PublicSessionID string
	RefreshToken    string
}

type RepoResetUserPasswordParam struct {
	Email           string
	ConfirmationKey string
}
