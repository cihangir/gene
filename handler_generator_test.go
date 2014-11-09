package gene

import (
	"encoding/json"
	"testing"
)

func TestHandlerGenerator(t *testing.T) {
	var s Schema
	if err := json.Unmarshal([]byte(testJSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	m := NewModule(&s)
	err := m.GenerateHandlers()
	if err != nil {
		t.Fatalf("err while generating handlers %s", err.Error())
	}
}
