package stringext

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var testData = []struct {
	Value                   string
	ToLowerFirst            string
	ToUpperFirst            string
	Pointerize              string
	JSONTag                 string
	JSONTagRequired         string
	Normalize               string
	ToFieldName             string
	DepunctWithInitialUpper string
	DepunctWithInitialLower string
}{
	{
		Value:                   "name",
		ToLowerFirst:            "name",
		ToUpperFirst:            "Name",
		Pointerize:              "n",
		JSONTag:                 "`json:\"name,omitempty\"`",
		JSONTagRequired:         "`json:\"name\"`",
		Normalize:               "name",
		ToFieldName:             "name",
		DepunctWithInitialUpper: "Name",
		DepunctWithInitialLower: "name",
	},
	{
		Value:                   "provider_id",
		ToLowerFirst:            "provider_id",
		ToUpperFirst:            "Provider_id",
		Pointerize:              "p",
		JSONTag:                 "`json:\"providerId,omitempty\"`",
		JSONTagRequired:         "`json:\"providerId\"`",
		Normalize:               "providerId",
		ToFieldName:             "provider_id",
		DepunctWithInitialUpper: "ProviderID",
		DepunctWithInitialLower: "providerID",
	},
	{
		Value:                   "app-identity",
		ToLowerFirst:            "app-identity",
		ToUpperFirst:            "App-identity",
		Pointerize:              "a",
		JSONTag:                 "`json:\"appIdentity,omitempty\"`",
		JSONTagRequired:         "`json:\"appIdentity\"`",
		Normalize:               "appIdentity",
		ToFieldName:             "app_identity",
		DepunctWithInitialUpper: "AppIdentity",
		DepunctWithInitialLower: "appIdentity",
	},
	{
		Value:                   "uuid",
		ToLowerFirst:            "uuid",
		ToUpperFirst:            "Uuid",
		Pointerize:              "u",
		JSONTag:                 "`json:\"uuid,omitempty\"`",
		JSONTagRequired:         "`json:\"uuid\"`",
		Normalize:               "uuid",
		ToFieldName:             "uuid",
		DepunctWithInitialUpper: "UUID",
		DepunctWithInitialLower: "uuid",
	},
	{
		Value:                   "oauth-client",
		ToLowerFirst:            "oauth-client",
		ToUpperFirst:            "Oauth-client",
		Pointerize:              "o",
		JSONTag:                 "`json:\"oauthClient,omitempty\"`",
		JSONTagRequired:         "`json:\"oauthClient\"`",
		Normalize:               "oauthClient",
		ToFieldName:             "oauth_client",
		DepunctWithInitialUpper: "OAuthClient",
		DepunctWithInitialLower: "oauthClient",
	},
	{
		Value:                   "Dyno all",
		ToLowerFirst:            "dyno all",
		ToUpperFirst:            "Dyno all",
		Pointerize:              "d",
		JSONTag:                 "`json:\"dynoAll,omitempty\"`",
		JSONTagRequired:         "`json:\"dynoAll\"`",
		Normalize:               "DynoAll",
		ToFieldName:             "dyno_all",
		DepunctWithInitialUpper: "DynoAll",
		DepunctWithInitialLower: "DynoAll",
	},
	{
		Value:                   "providerId",
		ToLowerFirst:            "providerId",
		ToUpperFirst:            "ProviderId",
		Pointerize:              "p",
		JSONTag:                 "`json:\"providerId,omitempty\"`",
		JSONTagRequired:         "`json:\"providerId\"`",
		Normalize:               "providerId",
		ToFieldName:             "provider_id",
		DepunctWithInitialUpper: "ProviderID",
		DepunctWithInitialLower: "providerID",
	},
	{
		Value:                   "Id",
		ToLowerFirst:            "id",
		ToUpperFirst:            "Id",
		Pointerize:              "i",
		JSONTag:                 "`json:\"id,omitempty\"`",
		JSONTagRequired:         "`json:\"id\"`",
		Normalize:               "Id",
		ToFieldName:             "id",
		DepunctWithInitialUpper: "ID",
		DepunctWithInitialLower: "ID",
	},
}

func TestToLowerFirst(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.ToLowerFirst, ToLowerFirst(ict.Value))
	}
}

func TestUpperFirst(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.ToUpperFirst, ToUpperFirst(ict.Value))
	}
}

func TestPointerize(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.Pointerize, Pointerize(ict.Value))
	}
}

func TestJSONTag(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.JSONTag, JSONTag(ict.Value, false))
	}
}

func TestJSONTagRequired(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.JSONTagRequired, JSONTag(ict.Value, true))
	}
}

func TestNormalize(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.Normalize, Normalize(ict.Value))
	}
}

func TestToFieldName(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.ToFieldName, ToFieldName(ict.Value))
	}
}

func TestDepunctWithInitialUpper(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.DepunctWithInitialUpper, DepunctWithInitialUpper(ict.Value))
	}
}

func TestDepunctWithInitialLower(t *testing.T) {
	for _, ict := range testData {
		equals(t, ict.DepunctWithInitialLower, DepunctWithInitialLower(ict.Value))
	}
}

// func TestInitialCap(t *testing.T) {
// 	for i, ict := range initialCapTests {
// 		depuncted := Depunct(ict.Ident, true)
// 		if depuncted != ict.Depuncted {
// 			t.Errorf("%d: wants %v, got %v", i, ict.Depuncted, depuncted)
// 		}
// 	}
// }

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}
