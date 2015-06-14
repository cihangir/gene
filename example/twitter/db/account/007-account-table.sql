
-- ----------------------------
--  Table structure for account.account
-- ----------------------------
DROP TABLE IF EXISTS "account"."account";
CREATE TABLE "account"."account" (
    "id" BIGINT DEFAULT nextval('account.account_id_seq' :: regclass)
        CONSTRAINT "check_account_id_gte_0" CHECK ("id" >= 0.000000),
    "profile_id" BIGINT
        CONSTRAINT "check_account_profile_id_gte_0" CHECK ("profile_id" >= 0.000000),
    "password" TEXT COLLATE "default"
        CONSTRAINT "check_account_password_min_length_6" CHECK (char_length("password") > 6 ),
    "url" TEXT COLLATE "default"
        CONSTRAINT "check_account_url_min_length_6" CHECK (char_length("url") > 6 ),
    "password_status_constant" "account"."account_password_status_constant_enum" DEFAULT 'valid',
    "salt" VARCHAR (255) COLLATE "default",
    "email_address" TEXT COLLATE "default",
    "email_status_constant" "account"."account_email_status_constant_enum" DEFAULT 'notVerified',
    "status_constant" "account"."account_status_constant_enum" DEFAULT 'registered',
    "created_at" TIMESTAMP (6) WITH TIME ZONE DEFAULT now()
) WITH (OIDS = FALSE);-- end schema creation
GRANT SELECT, INSERT, DELETE ON "account"."account" TO "twitter_db_role";