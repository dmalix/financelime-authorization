/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

DROP TABLE IF EXISTS "public"."account";
DROP SEQUENCE IF EXISTS "public"."account_id_seq";

DROP TABLE IF EXISTS "public"."session";
DROP SEQUENCE IF EXISTS "public"."session_id_seq";

DROP TABLE IF EXISTS "public"."device";
DROP SEQUENCE IF EXISTS "public"."device_id_seq";

DROP TABLE IF EXISTS "public"."invite_code";
DROP SEQUENCE IF EXISTS "public"."invite_code_id_seq";

DROP TABLE IF EXISTS "public"."invite_code_issued";
DROP SEQUENCE IF EXISTS "public"."invite_code_issued_id_seq";
