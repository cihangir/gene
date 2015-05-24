package definitions

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

func TestSchema(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.TestDataFull), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)

	context := config.NewContext()
	moduleName := context.ModuleNameFunc(s.Title)

	index := 0
	for _, def := range s.Definitions {

		// schema should have our generator
		if !def.Generators.Has(generatorName) {
			continue
		}

		settings, _ := def.Generators.Get(generatorName)
		settings.SetNX("schemaName", stringext.ToFieldName(moduleName))
		settings.SetNX("tableName", stringext.ToFieldName(def.Title))
		settings.SetNX("roleName", stringext.ToFieldName(moduleName))

		sts := DefineSchema(settings, def)
		equals(t, expectedSchemas[index], sts)
		index++
	}
}

var expectedSchemas = []string{
	`
-- ----------------------------
--  Schema structure for account
-- ----------------------------
-- create schema
CREATE SCHEMA IF NOT EXISTS "account";
-- give usage permission
GRANT usage ON SCHEMA "account" to "social";
-- add new schema to search path -just for convenience
-- SELECT set_config('search_path', current_setting('search_path') || ',account', false);`,
}
