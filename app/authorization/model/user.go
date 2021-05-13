package model

import "time"

type User struct {
	ID       int64
	Email    string
	Password string
	Language string
}

type Device struct {
	Platform string `json:"platform" example:"Linux x86_64"`
	Height   int    `json:"height" example:"1920"`
	Width    int    `json:"width" example:"1060"`
	Language string `json:"language" example:"en-US"`
	Timezone string `json:"timezone" example:"2"`
}

type Session struct {
	PublicSessionID string    `json:"sessionID"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Platform        string    `json:"platform"`
}

type InviteCodeRecord struct {
	Id          int64
	UserID      int64
	LimitAmount int
	value       string
}
