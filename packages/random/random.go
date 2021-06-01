/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package random

import (
	"math/rand"
)

func StringRand(min, max int, readable bool) string {

	const charSetReadable = "abcefghijkmnopqrtuvwxyz23479"
	const charSetUnreadable = "abcdefghijklmnopqrstuvwxyz0123456789"

	var length int
	var chars string

	if min < max {
		length = min + rand.Intn(max-min)
	} else {
		length = min
	}

	if readable {
		chars = charSetReadable
	} else {
		chars = charSetUnreadable
	}

	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = chars[rand.Intn(len(chars)-1)]
	}
	return string(buf)
}
