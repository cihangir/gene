package definitions

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

func TestSequence(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.TestDataFull), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)
	g := New()

	context := config.NewContext()
	moduleName := context.ModuleNameFunc(s.Title)
	settings := g.generateSettings(moduleName, s)

	index := 0
	for _, def := range s.Definitions {

		// schema should have our generator
		if !def.Generators.Has(generatorName) {
			continue
		}

		settingsDef := g.setDefaultSettings(settings, def)
		settingsDef.Set("tableName", stringext.ToFieldName(def.Title))

		sts := DefineSequence(settingsDef, def)
		equals(t, expectedSequences[index], sts)
		index++
	}
}

var expectedSequences = []string{
	`
-- ----------------------------
--  Sequence structure for account.profile_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "account"."profile_id_seq" CASCADE;
CREATE SEQUENCE "account"."profile_id_seq" INCREMENT 1 START 1 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
GRANT USAGE ON SEQUENCE "account"."profile_id_seq" TO "social";`,
}
