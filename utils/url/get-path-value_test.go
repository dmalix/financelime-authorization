/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package url

import (
	"testing"
)

func TestGetPathValue_Success(t *testing.T) {

	result, err:= GetPathValue("/test/target",1)

	if err != nil {
		t.Errorf("function returned wrong the err value: got %v want %v",
			err, nil)
	}

	if result != "target" {
		t.Errorf("function returned wrong the result value: got %v want %v",
			result, "target")
	}
}

func TestGetPathValue_Error(t *testing.T) {

	_, err:= GetPathValue("/test/target",2)

	if err == nil {
		t.Errorf("function returned wrong the err value: got %v want %v",
			err, nil)
	}
}
