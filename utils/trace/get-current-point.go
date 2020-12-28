/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package trace

import (
	"fmt"
	"runtime"
	"strings"
)

func GetCurrentPoint() string {

	const postfix = "-repo"
	const mod = "github.com/dmalix/financelime-authorization"
	var location string

	pc := make([]uintptr, 15)
	callers := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:callers])
	frame, more := frames.Next()
	file := strings.Split(frame.File, postfix)
	function := strings.Split(frame.Function, mod)
	fmt.Println(function)

	if len(file) != 0 && len(function) != 0 && more != false {
		location = fmt.Sprintf("%s:%d %s", file[1], frame.Line, function[len(function)-1])
	}
	return location
}
