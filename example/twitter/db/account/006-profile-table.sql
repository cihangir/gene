
-- ----------------------------
--  Table structure for account.profile
-- ----------------------------
DROP TABLE IF EXISTS "account"."profile";
CREATE TABLE "account"."profile" (
    "avatar_url" VARCHAR (2000) COLLATE "default",
    "created_at" TIMESTAMP (6) WITH TIME ZONE DEFAULT now(),
    "description" VARCHAR (160) COLLATE "default",
    "id" BIGINT DEFAULT nextval('account.profile_id_seq' :: regclass)
        CONSTRAINT "check_profile_id_gte_0" CHECK ("id" >= 0.000000),
    "link_color" VARCHAR (6) COLLATE "default" DEFAULT 'FF0000',
    "location" VARCHAR (30) COLLATE "default",
    "screen_name" VARCHAR (20) COLLATE "default"
        CONSTRAINT "check_profile_screen_name_min_length_4" CHECK (char_length("screen_name") > 4 ),
    "url" VARCHAR (100) COLLATE "default"
) WITH (OIDS = FALSE);-- end schema creation
GRANT SELECT, INSERT, DELETE ON "account"."profile" TO "twitter_db_role";