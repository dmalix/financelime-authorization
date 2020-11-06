/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

//noinspection GoSnakeCaseUsage
type Dist_MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var Dist_MockData Dist_MockDescription

func (a *Dist_MockDescription) Dist() (string, string, error) {
	return "version", "build", Dist_MockData.Expected.Error
}
