/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"context"
	"go.uber.org/zap"
	"testing"
)

func TestDist(t *testing.T) {

	var ctx context.Context
	newLogger := new(zap.Logger)
	newService := NewService("version", "buildTime", "commit", "compiler")

	_, _, err := newService.Version(ctx, newLogger)
	if err != nil {
		t.Errorf("Service returned wrong err value: got %v want %v",
			err, nil)
	}
}
