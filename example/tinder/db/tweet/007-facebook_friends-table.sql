-- ----------------------------
--  Table structure for tweet.facebook_friends
-- ----------------------------
DROP TABLE IF EXISTS "tweet"."facebook_friends";
CREATE TABLE "tweet"."facebook_friends" (
    "source_id" TEXT COLLATE "default"
        CONSTRAINT "check_facebook_friends_source_id_min_length_1" CHECK (char_length("source_id") > 1 ),
    "target_id" TEXT COLLATE "default"
        CONSTRAINT "check_facebook_friends_target_id_min_length_1" CHECK (char_length("target_id") > 1 )
) WITH (OIDS = FALSE);-- end schema creation
GRANT SELECT, INSERT, DELETE ON "tweet"."facebook_friends" TO "twitter_db_role";