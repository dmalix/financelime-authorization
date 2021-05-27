/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

type service struct {
	VersionNumber    string
	VersionBuildTime string
	VersionCommit    string
	VersionCompiler  string
}

func NewService(
	versionNumber string,
	versionBuildTime string,
	versionCommit string,
	versionCompiler string) *service {
	return &service{
		VersionNumber:    versionNumber,
		VersionBuildTime: versionBuildTime,
		VersionCommit:    versionCommit,
		VersionCompiler:  versionCompiler,
	}
}

func (service *service) Version(_ context.Context, _ *zap.Logger) (string, string, error) {

	versionNumber := service.VersionNumber
	versionBuild := fmt.Sprintf("%s [%s]", service.VersionCommit, service.VersionBuildTime)

	return versionNumber, versionBuild, nil
}
