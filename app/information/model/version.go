package model

type VersionResponse struct {
	Number string `json:"number" example:"v0.2.0-alpha"`
	Build  string `json:"build" example:"fc56bb1 [2021-05-07_11:12:09]"`
}
