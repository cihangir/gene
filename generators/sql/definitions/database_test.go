package definitions

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

func TestDatabase(t *testing.T) {
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

		sts, err := DefineDatabase(settingsDef, def)
		if err != nil {
			t.Fatal(err.Error())
		}

		equals(t, expectedDatabases[index], string(sts))
		index++
	}
}

var expectedDatabases = []string{
	`
-- Drop database
DROP DATABASE IF EXISTS "mydatabase";
-- Create database itself
CREATE DATABASE "mydatabase" OWNER "social" ENCODING 'UTF8'  TEMPLATE template0;`,
}
