/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

//noinspection GoSnakeCaseUsage
type MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var MockData MockDescription

func (a *MockDescription) Dist() (string, string, error) {
	return "version", "build", MockData.Expected.Error
}
