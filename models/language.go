/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package models

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
