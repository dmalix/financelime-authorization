/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package trace

type MockDescription struct {
	Props    struct{}
	Expected struct {
		Error error
	}
}

func (s *MockDescription) GetCurrentLocation() string {
	return "value"
}
