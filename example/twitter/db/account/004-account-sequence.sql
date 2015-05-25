
-- ----------------------------
--  Sequence structure for account.account_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "account"."account_id_seq" CASCADE;
CREATE SEQUENCE "account"."account_id_seq" INCREMENT 1 START 1 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
GRANT USAGE ON SEQUENCE "account"."account_id_seq" TO "twitter_db_role";