package gene

import (
	"encoding/json"
	"testing"
)

func TestCreateModuleStructure(t *testing.T) {
	expected := []string{
		"gene/modules/name",
		"gene/modules/name/api",
		"gene/modules/name/name",
		"gene/modules/name/cmd",
		"gene/modules/name/tests",
		"gene/modules/name/errors",
		"gene/modules/name/handlers",
	}

	structure := createModuleStructure("name")

	for _, stc := range structure {
		exists := false
		for _, expt := range expected {
			if expt == stc {
				exists = true
				break
			}
		}
		if !exists {
			t.Fatalf("%s is not expected in the result set", stc)
		}
	}
}

func TestCreateModule(t *testing.T) {
	var s Schema
	if err := json.Unmarshal([]byte(testJSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	m := NewModule(&s)
	m.Create()
}
