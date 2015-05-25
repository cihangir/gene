package modules

import "testing"

func TestCreateModuleStructure(t *testing.T) {
	expected := []string{
		"cmd/name/",

		"workers/name",
		"workers/name/api",
		"workers/name/tests",
		"workers/name/js",
		"workers/name/errors",
		"workers/name/handlers",
		"workers/name/clients",
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
