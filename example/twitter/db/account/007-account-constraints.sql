-------------------------------
--  Primary key structure for table account
-- ----------------------------
ALTER TABLE "account"."account" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;
-------------------------------
--  Uniqueu key structure for table account
-- ----------------------------
ALTER TABLE "account"."account" ADD CONSTRAINT "key_account_id" UNIQUE ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "account"."account" ADD CONSTRAINT "key_account_url" UNIQUE ("url") NOT DEFERRABLE INITIALLY IMMEDIATE;