/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

package authorization

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

func (repo *RepoMockDescription) createUser(_ repoCreateUserParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) confirmUserEmail(_ string) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) getUserByAuth(_ repoGetUserByAuthParam) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) saveSession(_ repoSaveSessionParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) getListActiveSessions(userID int64) ([]session, error) {
	var sessions []session
	return sessions, RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) updateSession(_ repoUpdateSessionParam) error {
	return RepoMockData.Expected.Error
}

func (repo *RepoMockDescription) getUserByRefreshToken(RefreshToken string) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}

func (s *RepoMockDescription) deleteSession(_ repoDeleteSessionParam) error {
	return RepoMockData.Expected.Error
}

func (s *RepoMockDescription) requestUserPasswordReset(_ repoRequestUserPasswordResetParam) (user, error) {
	var user user
	return user, RepoMockData.Expected.Error
}
