package authorization

import "time"

type inviteCodeRecord struct {
	id          int64
	userID      int64
	numberLimit int
	value       string
}

type device struct {
	Platform string `json:"platform" example:"Linux x86_64"`
	Height   int    `json:"height" example:"1920"`
	Width    int    `json:"width" example:"1060"`
	Language string `json:"language" example:"en-US"`
	Timezone string `json:"timezone" example:"2"`
}

type user struct {
	id       int64
	email    string
	password string
	language string
}

type session struct {
	PublicSessionID string    `json:"sessionID"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Platform        string    `json:"platform"`
}

type LanguageContent struct {
	Language map[string]int
	Data     DataLanguageContent
}

type FileLanguageContent struct {
	Language []string
	Data     DataLanguageContent
}

type DataLanguageContent struct {
	User UserDataLanguageContent
}

type UserDataLanguageContent struct {
	Signup struct {
		Email struct {
			Confirm struct {
				Subject []string
				Body    []string
			}
			Password struct {
				Subject []string
				Body    []string
			}
		}
		Page struct {
			Text []string
		}
	}
	ResetPassword struct {
		Email struct {
			Request struct {
				Subject []string
				Body    []string
			}
			Password struct {
				Subject []string
				Body    []string
			}
		}
		Page struct {
			Request struct {
				Text []string
			}
		}
	}
	Login struct {
		Email struct {
			Subject []string
			Body    []string
		}
	}
}
