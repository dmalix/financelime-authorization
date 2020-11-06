/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import "fmt"

func (service *Service) Dist() (string, string, error) {

	version := service.version
	build := fmt.Sprintf("%s [%s]", service.commit, service.buildTime)

	return version, build, nil
}
