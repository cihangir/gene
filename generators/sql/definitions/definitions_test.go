package definitions

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"testing"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestDefinitions(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.TestDataFull), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)

	sts, err := New().Generate(config.NewContext(), s)
	equals(t, nil, err)
	for _, s := range sts {
		equals(t, expected, string(s.Content))
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}

const expected = `
-- ----------------------------
--  Schema structure for account
-- ----------------------------
-- create schema
CREATE SCHEMA IF NOT EXISTS "account";
-- give usage permission
GRANT usage ON SCHEMA "account" to "social";
-- add new schema to search path -just for convenience
-- SELECT set_config('search_path', current_setting('search_path') || ',account', false);
-- ----------------------------
--  Sequence structure for account.profile_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "account"."profile_id_seq" CASCADE;
CREATE SEQUENCE "account"."profile_id_seq" INCREMENT 1 START 1 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
GRANT USAGE ON SEQUENCE "account"."profile_id_seq" TO "social";
-- ----------------------------
--  Required extensions
-- ----------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- ----------------------------
--  Types structure for account.profile.enum_bare
-- ----------------------------
DROP TYPE IF EXISTS "account"."profile_enum_bare_enum" CASCADE;
CREATE TYPE "account"."profile_enum_bare_enum" AS ENUM (
  'enum1', 'enum2', 'enum3'
);
-- ----------------------------
--  Table structure for account.profile
-- ----------------------------
DROP TABLE IF EXISTS "account"."profile";
CREATE TABLE "account"."profile" (
    "boolean_bare" BOOLEAN,
    "boolean_with_default" BOOLEAN DEFAULT TRUE,
    "boolean_with_max_length" BOOLEAN,
    "boolean_with_min_length" BOOLEAN,
    "enum_bare" "profile_enum_bare_enum",
    "id" BIGINT DEFAULT nextval('account.profile_id_seq' :: regclass)
        CONSTRAINT "check_profile_id_gte_0" CHECK ("id" >= 0.000000),
    "number_bare" NUMERIC,
    "number_with_exclusive_maximum_without_maximum" NUMERIC,
    "number_with_exclusive_minimum" NUMERIC
        CONSTRAINT "check_profile_number_with_exclusive_minimum_gte_0" CHECK ("number_with_exclusive_minimum" >= 0.000000),
    "number_with_exclusive_minimum_without_minimum" NUMERIC,
    "number_with_maximum" NUMERIC
        CONSTRAINT "check_profile_number_with_maximum_lte_1023" CHECK ("number_with_maximum" <= 1023.000000),
    "number_with_maximum_as_float32" NUMERIC
        CONSTRAINT "check_profile_number_with_maximum_as_float32_lte_3" CHECK ("number_with_maximum_as_float32" <= 3.200000),
    "number_with_maximum_as_float64" NUMERIC
        CONSTRAINT "check_profile_number_with_maximum_as_float64_lte_6" CHECK ("number_with_maximum_as_float64" <= 6.400000),
    "number_with_maximum_as_int" INTEGER
        CONSTRAINT "check_profile_number_with_maximum_as_int_lte_2" CHECK ("number_with_maximum_as_int" <= 2.000000),
    "number_with_maximum_as_int16" SMALLINT
        CONSTRAINT "check_profile_number_with_maximum_as_int16_lte_2" CHECK ("number_with_maximum_as_int16" <= 2.000000),
    "number_with_maximum_as_int32" INTEGER
        CONSTRAINT "check_profile_number_with_maximum_as_int32_lte_2" CHECK ("number_with_maximum_as_int32" <= 2.000000),
    "number_with_maximum_as_int64" BIGINT
        CONSTRAINT "check_profile_number_with_maximum_as_int64_lte_64" CHECK ("number_with_maximum_as_int64" <= 64.000000),
    "number_with_maximum_as_int8" SMALLINT
        CONSTRAINT "check_profile_number_with_maximum_as_int8_lte_2" CHECK ("number_with_maximum_as_int8" <= 2.000000),
    "number_with_maximum_as_u_int" INTEGER
        CONSTRAINT "check_profile_number_with_maximum_as_u_int_lte_2" CHECK ("number_with_maximum_as_u_int" <= 2.000000),
    "number_with_maximum_as_u_int16" SMALLINT
        CONSTRAINT "check_profile_number_with_maximum_as_u_int16_lte_2" CHECK ("number_with_maximum_as_u_int16" <= 2.000000),
    "number_with_maximum_as_u_int32" INTEGER
        CONSTRAINT "check_profile_number_with_maximum_as_u_int32_lte_2" CHECK ("number_with_maximum_as_u_int32" <= 2.000000),
    "number_with_maximum_as_u_int64" BIGINT
        CONSTRAINT "check_profile_number_with_maximum_as_u_int64_lte_64" CHECK ("number_with_maximum_as_u_int64" <= 64.000000),
    "number_with_maximum_as_u_int8" SMALLINT
        CONSTRAINT "check_profile_number_with_maximum_as_u_int8_lte_2" CHECK ("number_with_maximum_as_u_int8" <= 2.000000),
    "number_with_minimum_as_float32" NUMERIC
        CONSTRAINT "check_profile_number_with_minimum_as_float32_gte_0" CHECK ("number_with_minimum_as_float32" >= 0.000000),
    "number_with_minimum_as_float64" NUMERIC
        CONSTRAINT "check_profile_number_with_minimum_as_float64_gte_0" CHECK ("number_with_minimum_as_float64" >= 0.000000),
    "number_with_minimum_as_int" INTEGER
        CONSTRAINT "check_profile_number_with_minimum_as_int_gte_0" CHECK ("number_with_minimum_as_int" >= 0.000000),
    "number_with_minimum_as_int16" SMALLINT
        CONSTRAINT "check_profile_number_with_minimum_as_int16_gte_0" CHECK ("number_with_minimum_as_int16" >= 0.000000),
    "number_with_minimum_as_int32" INTEGER
        CONSTRAINT "check_profile_number_with_minimum_as_int32_gte_0" CHECK ("number_with_minimum_as_int32" >= 0.000000),
    "number_with_minimum_as_int64" BIGINT
        CONSTRAINT "check_profile_number_with_minimum_as_int64_gte_0" CHECK ("number_with_minimum_as_int64" >= 0.000000),
    "number_with_minimum_as_int8" SMALLINT
        CONSTRAINT "check_profile_number_with_minimum_as_int8_gte_0" CHECK ("number_with_minimum_as_int8" >= 0.000000),
    "number_with_minimum_as_u_int" INTEGER
        CONSTRAINT "check_profile_number_with_minimum_as_u_int_gte_0" CHECK ("number_with_minimum_as_u_int" >= 0.000000),
    "number_with_minimum_as_u_int16" SMALLINT
        CONSTRAINT "check_profile_number_with_minimum_as_u_int16_gte_0" CHECK ("number_with_minimum_as_u_int16" >= 0.000000),
    "number_with_minimum_as_u_int32" INTEGER
        CONSTRAINT "check_profile_number_with_minimum_as_u_int32_gte_0" CHECK ("number_with_minimum_as_u_int32" >= 0.000000),
    "number_with_minimum_as_u_int64" BIGINT
        CONSTRAINT "check_profile_number_with_minimum_as_u_int64_gte_0" CHECK ("number_with_minimum_as_u_int64" >= 0.000000),
    "number_with_minimum_as_u_int8" SMALLINT
        CONSTRAINT "check_profile_number_with_minimum_as_u_int8_gte_0" CHECK ("number_with_minimum_as_u_int8" >= 0.000000),
    "number_with_multiple_of" NUMERIC
        CONSTRAINT "check_profile_number_with_multiple_of_multiple_of_2" CHECK (("number_with_multiple_of" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_float32" NUMERIC
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_float32_multiple_of_3" CHECK (("number_with_multiple_of_formatted_as_float32" % 3.200000) = 0),
    "number_with_multiple_of_formatted_as_float64" NUMERIC
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_float64_multiple_of_6" CHECK (("number_with_multiple_of_formatted_as_float64" % 6.400000) = 0),
    "number_with_multiple_of_formatted_as_int" INTEGER
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_int_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_int" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_int16" SMALLINT
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_int16_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_int16" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_int32" INTEGER
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_int32_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_int32" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_int64" BIGINT
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_int64_multiple_of_64" CHECK (("number_with_multiple_of_formatted_as_int64" % 64.000000) = 0),
    "number_with_multiple_of_formatted_as_int8" SMALLINT
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_int8_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_int8" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_u_int" INTEGER
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_u_int_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_u_int" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_u_int16" SMALLINT
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_u_int16_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_u_int16" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_u_int32" INTEGER
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_u_int32_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_u_int32" % 2.000000) = 0),
    "number_with_multiple_of_formatted_as_u_int64" BIGINT
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_u_int64_multiple_of_64" CHECK (("number_with_multiple_of_formatted_as_u_int64" % 64.000000) = 0),
    "number_with_multiple_of_formatted_as_u_int8" SMALLINT
        CONSTRAINT "check_profile_number_with_multiple_of_formatted_as_u_int8_multiple_of_2" CHECK (("number_with_multiple_of_formatted_as_u_int8" % 2.000000) = 0),
    "string_bare" TEXT COLLATE "default",
    "string_date_formatted" TIMESTAMP (6) WITH TIME ZONE,
    "string_date_formatted_with_default" TIMESTAMP (6) WITH TIME ZONE DEFAULT now(),
    "string_uuid_formatted" UUID,
    "string_uuid_formatted_with_default" UUID DEFAULT uuid_generate_v1(),
    "string_with_default" TEXT COLLATE "default" DEFAULT 'THISISMYDEFAULTVALUE',
    "string_with_max_and_min_length" VARCHAR (24) COLLATE "default"
        CONSTRAINT "check_profile_string_with_max_and_min_length_min_length_4" CHECK (char_length("string_with_max_and_min_length") > 4 ),
    "string_with_max_length" VARCHAR (24) COLLATE "default",
    "string_with_min_length" TEXT COLLATE "default"
        CONSTRAINT "check_profile_string_with_min_length_min_length_24" CHECK (char_length("string_with_min_length") > 24 ),
    "string_with_pattern" TEXT COLLATE "default"
        CONSTRAINT "check_profile_string_with_pattern_pattern" CHECK ("string_with_pattern" ~ '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
) WITH (OIDS = FALSE);-- end schema creation
GRANT SELECT, UPDATE ON "account"."profile" TO "social";`
