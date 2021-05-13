/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

import "context"

//noinspection GoSnakeCaseUsage
type RepoMockDescription struct {
	Props struct {
	}
	Expected struct {
		Error error
	}
}

//noinspection GoSnakeCaseUsage
var RepoMockData RepoMockDescription

func (repo *RepoMockDescription) createUser(_ context.Context, _ repoCreateUserParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) confirmUserEmail(_ context.Context, _ string) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) getUserByAuth(_ context.Context, _ repoGetUserByAuthParam) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) saveSession(_ context.Context, _ repoSaveSessionParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) getListActiveSessions(_ context.Context, _ int64) ([]session, error) {
	var sessions []session
	return sessions, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) updateSession(_ context.Context, _ repoUpdateSessionParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) getUserByRefreshToken(_ context.Context, _ string) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) deleteSession(_ context.Context, _ repoDeleteSessionParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) requestUserPasswordReset(_ context.Context, _ repoRequestUserPasswordResetParam) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}
