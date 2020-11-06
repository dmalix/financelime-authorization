/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

type Service struct {
	version   string
	buildTime string
	commit    string
	compiler  string
}

func NewService(
	version string,
	buildTime string,
	commit string,
	compiler string) *Service {
	return &Service{
		version:   version,
		buildTime: buildTime,
		commit:    commit,
		compiler:  compiler,
	}
}
