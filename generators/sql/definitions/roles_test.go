package definitions

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

func TestRoles(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.TestDataFull), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)
	g := New()

	context := common.NewContext()
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

		sts, err := DefineRole(context, settingsDef, def)
		if err != nil {
			t.Fatal(err.Error())
		}

		equals(t, expectedRoles[index], string(sts))
		index++
	}
}

var expectedRoles = []string{
	`
-- Drop role
DROP ROLE IF EXISTS "social";
-- Create role
CREATE ROLE "social";
-- Create shadow user for future extensibility
DROP USER IF EXISTS "socialapplication";
CREATE USER "socialapplication" PASSWORD 'socialapplication';
GRANT "social" TO "socialapplication";`,
}
