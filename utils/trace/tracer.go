/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package trace

type GetCurrentLocation interface {
	GetCurrentLocation() string
}

type Config struct {
	Mod string
}

type Tracer struct {
	config Config
}

func NewTracer(
	config Config) *Tracer {
	return &Tracer{
		config: config,
	}
}