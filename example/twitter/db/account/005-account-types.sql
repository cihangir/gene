
-- ----------------------------
--  Types structure for account.account.email_status_constant
-- ----------------------------
DROP TYPE IF EXISTS "account"."account_email_status_constant_enum" CASCADE;
CREATE TYPE "account"."account_email_status_constant_enum" AS ENUM (
  'verified',
  'notVerified'
);
ALTER TYPE "account"."account_email_status_constant_enum" OWNER TO "twitter_db_role";
-- ----------------------------
--  Types structure for account.account.password_status_constant
-- ----------------------------
DROP TYPE IF EXISTS "account"."account_password_status_constant_enum" CASCADE;
CREATE TYPE "account"."account_password_status_constant_enum" AS ENUM (
  'valid',
  'needsReset',
  'generated'
);
ALTER TYPE "account"."account_password_status_constant_enum" OWNER TO "twitter_db_role";
-- ----------------------------
--  Types structure for account.account.status_constant
-- ----------------------------
DROP TYPE IF EXISTS "account"."account_status_constant_enum" CASCADE;
CREATE TYPE "account"."account_status_constant_enum" AS ENUM (
  'registered',
  'unregistered',
  'needsManualVerification'
);
ALTER TYPE "account"."account_status_constant_enum" OWNER TO "twitter_db_role";