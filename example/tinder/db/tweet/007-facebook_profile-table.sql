-- ----------------------------
--  Table structure for tweet.facebook_profile
-- ----------------------------
DROP TABLE IF EXISTS "tweet"."facebook_profile";
CREATE TABLE "tweet"."facebook_profile" (
    "id" TEXT COLLATE "default" DEFAULT nextval('tweet.facebook_profile_id_seq' :: regclass),
    "first_name" TEXT COLLATE "default"
        CONSTRAINT "check_facebook_profile_first_name_min_length_1" CHECK (char_length("first_name") > 1 ),
    "middle_name" TEXT COLLATE "default",
    "last_name" TEXT COLLATE "default",
    "picture_url" TEXT COLLATE "default"
) WITH (OIDS = FALSE);-- end schema creation
GRANT SELECT, INSERT, UPDATE ON "tweet"."facebook_profile" TO "twitter_db_role";