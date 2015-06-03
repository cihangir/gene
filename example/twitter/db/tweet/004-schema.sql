
-- ----------------------------
--  Schema structure for tweet
-- ----------------------------
-- create schema
CREATE SCHEMA IF NOT EXISTS "tweet";
-- give usage permission
GRANT usage ON SCHEMA "tweet" to "twitter_db_role";
-- add new schema to search path -just for convenience
-- SELECT set_config('search_path', current_setting('search_path') || ',tweet', false);