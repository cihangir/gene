
-- ----------------------------
--  Table structure for tweet.tweet
-- ----------------------------
DROP TABLE IF EXISTS "tweet"."tweet";
CREATE TABLE "tweet"."tweet" (
    "id" BIGINT DEFAULT nextval('tweet.tweet_id_seq' :: regclass),
    "profile_id" BIGINT
        CONSTRAINT "check_tweet_profile_id_gte_0" CHECK ("profile_id" >= 0.000000),
    "body" TEXT COLLATE "default"
        CONSTRAINT "check_tweet_body_min_length_1" CHECK (char_length("body") > 1 ),
    "location" TEXT COLLATE "default",
    "retweet_count" INTEGER,
    "favourites_count" INTEGER,
    "possibly_sensitive" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP (6) WITH TIME ZONE DEFAULT now()
) WITH (OIDS = FALSE);-- end schema creation
GRANT SELECT, INSERT, DELETE ON "tweet"."tweet" TO "twitter_db_role";