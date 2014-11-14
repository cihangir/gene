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
	if err := m.Create(); err != nil {

		t.Fatalf("err while generating handlers %s", err.Error())
	}

	// if err := m.GenerateHandlers(); err != nil {
	// 	t.Fatalf("err while generating handlers %s", err.Error())
	// }
}
