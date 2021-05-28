/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package email

//noinspection GoSnakeCaseUsage
type AddEmailMessageToQueue_MockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var AddEmailMessageToQueue_MockData AddEmailMessageToQueue_MockDescription

func (manager AddEmailMessageToQueue_MockDescription) AddEmailMessageToQueue(_ chan MessageBox, _ Request, _ Email) error {
	return AddEmailMessageToQueue_MockData.Expected.Error
}
