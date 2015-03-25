package constants

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"testing"

	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

const expected = `
	// AccountEmailStatusConstant holds the predefined enums
	var AccountEmailStatusConstant = struct {
		Verified    string
		NotVerified string
	}{
		Verified:    "verified",
		NotVerified: "notVerified",
	}

	// AccountPasswordStatusConstant holds the predefined enums
	var AccountPasswordStatusConstant = struct {
		Valid      string
		NeedsReset string
		Generated  string
	}{
		Valid:      "valid",
		NeedsReset: "needsReset",
		Generated:  "generated",
	}

	// AccountStatusConstant holds the predefined enums
	var AccountStatusConstant = struct {
		Registered              string
		Unregistered            string
		NeedsManualVerification string
	}{
		Registered:              "registered",
		Unregistered:            "unregistered",
		NeedsManualVerification: "needsManualVerification",
	}`

func TestConstants(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(nil)

	a, err := Generate(s.Definitions["Account"])
	equals(t, nil, err)
	equals(t, expected, string(a))
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}
