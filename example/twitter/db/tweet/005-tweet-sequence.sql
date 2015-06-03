
-- ----------------------------
--  Sequence structure for tweet.tweet_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "tweet"."tweet_id_seq" CASCADE;
CREATE SEQUENCE "tweet"."tweet_id_seq" INCREMENT 1 START 1 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
GRANT USAGE ON SEQUENCE "tweet"."tweet_id_seq" TO "twitter_db_role";