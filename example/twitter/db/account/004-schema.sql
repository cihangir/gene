-- ----------------------------
--  Schema structure for account
-- ----------------------------
CREATE SCHEMA IF NOT EXISTS "account";
--
-- Give usage permission
--
GRANT usage ON SCHEMA "account" to "twitter_db_role";
--
-- add new schema to search path -just for convenience
-- SELECT set_config('search_path', current_setting('search_path') || ',account', false);