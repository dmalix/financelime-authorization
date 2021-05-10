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
	var fileName string
	var functionName string

	pc := make([]uintptr, 15)
	callers := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:callers])
	frame, _ := frames.Next()

	if strings.Contains(frame.File, postfix) {
		file := strings.Split(frame.File, postfix)
		fileName = file[len(file)-1]
	} else {
		fileName = frame.File
	}

	if strings.Contains(frame.Function, mod) {
		function := strings.Split(frame.Function, mod)
		functionName = function[len(function)-1]
	} else {
		functionName = frame.Function
	}

	location = fmt.Sprintf("%s:%d %s", fileName, frame.Line, functionName)

	return location
}
