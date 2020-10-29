/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

BEGIN;
	LOCK TABLE "user", invite_code IN SHARE MODE;
	INSERT INTO "user" ( created_at, email, "language", "password" )
	VALUES
		( NOW( ), 'dmalix@financelime.com', 'en', '594ceaa7946ce360a86f727a2d2c78162e67174a0044bfd143013c2921d2c048' );
	INSERT INTO invite_code ( created_at, user_id, expires_at, number_limit, "value" )
	VALUES
		( NOW( ), 1, '2021-12-31 21:23:21', 2, 'testInviteCode' );
COMMIT;

