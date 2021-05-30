package config

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
			Text []string
		}
	}
	Login struct {
		Email struct {
			Subject []string
			Body    []string
		}
	}
}
