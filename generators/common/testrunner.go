package common

import (
	"encoding/json"
	"testing"

	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func RunTest(t *testing.T, g Generator, expecteds []string) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)

	req := &Req{
		Schema:  s,
		Context: NewContext(),
	}

	res := &Res{}
	err := g.Generate(req, res)
	if err != nil {
		t.Fatal(err.Error())
	}

	for i, s := range res.Output {
		TestEquals(t, expecteds[i], string(s.Content))
	}
}
