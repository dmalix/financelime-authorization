/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package generator

import (
	"strconv"
	"testing"
)

func TestStringRandTrueMode(t *testing.T) {

	var values []string
	var value string
	var sum int

	for i := 0; i < 10000000; i++ {
		value = StringRand(16, 16, true)
		for j := 0; j < len(values); j++ {
			if values[j] == value {
				sum = sum + 1
			}
			values = append(values, value)
		}
	}

	if sum != 0 {
		t.Errorf("Duplicates found: %s", strconv.Itoa(sum))
	}
}
