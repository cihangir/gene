package constructors

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestConstructors(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	a, err := Generate(&s)
	common.TetEquals(t, nil, err)
	common.TetEquals(t, expected, string(a))
}

const expected = `
// NewAccount creates a new Account struct with default values
func NewAccount() *Account {
	return &Account{}
}`
