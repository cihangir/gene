-------------------------------
--  Primary key structure for table tweet
-- ----------------------------
ALTER TABLE "tweet"."tweet" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;
-------------------------------
--  Unique key structure for table tweet
-- ----------------------------
ALTER TABLE "tweet"."tweet" ADD CONSTRAINT "key_tweet_id" UNIQUE ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "tweet"."tweet" ADD CONSTRAINT "key_tweet_profile_id_body" UNIQUE ("profile_id", "body") NOT DEFERRABLE INITIALLY IMMEDIATE;