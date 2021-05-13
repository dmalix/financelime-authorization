/* Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
   Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
   License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html */

DROP TABLE IF EXISTS "public"."confirmation_create_new_user";
DROP SEQUENCE IF EXISTS "public"."confirmation_create_new_user_id_seq";

DROP TABLE IF EXISTS "public"."confirmation_reset_password";
DROP SEQUENCE IF EXISTS "public"."confirmation_reset_password_id_seq";

DROP TABLE IF	EXISTS "public"."invite_code_reserved";
DROP SEQUENCE IF EXISTS "public"."invite_code_reserved_id_seq";

DROP TABLE IF	EXISTS "public"."notification_email";
DROP SEQUENCE IF EXISTS "public"."notification_email_id_seq";

