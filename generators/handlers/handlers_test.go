package handlers

import (
	"encoding/json"
	"testing"

	"bitbucket.org/cihangirsavas/gene/generators/modules"
	"bitbucket.org/cihangirsavas/gene/schema"
	"bitbucket.org/cihangirsavas/gene/testdata"
)

func TestHandlerGenerator(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	m := modules.NewModule(&s)
	if err := m.Create(); err != nil {

		t.Fatalf("err while generating handlers %s", err.Error())
	}

	// if err := m.GenerateHandlers(); err != nil {
	// 	t.Fatalf("err while generating handlers %s", err.Error())
	// }
}
