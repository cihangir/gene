package schema

import (
	"encoding/json"

	"strings"
	"testing"

	"bitbucket.org/cihangirsavas/gene/testdata"
)

func TestGenerateSchema(t *testing.T) {
	var s Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	// replace "~" with "`"
	result := strings.Replace(`
// Message represents a simple post
type Message struct {
	Age            int       ~json:"age"~
	Body           string    ~json:"body"~ // The body for a message
	CreatedAt      time.Time ~json:"createdAt"~
	Enabled        bool      ~json:"enabled"~
	ID             int64     ~json:"id"~ // The unique identifier for a message
	StatusConstant string    ~json:"statusConstant"~
	Token          string    ~json:"token"~ // The token for a message security
}`, "~", "`", -1)

	code, err := s.GenerateSchema()
	if err != nil {
		t.Fatal(err.Error())
	}

	if result != string(code) {
		// fmt.Printf("foo %# v", pretty.Formatter(difflib.Diff([]string{result}, []string{string(code)})))
		t.Fatalf("Schema is not same. Wanted: %s, Get: %s", result, string(code))
	}
}

func TestGenerateValidators(t *testing.T) {
	var s Schema
	if err := json.Unmarshal([]byte(testJSON1), &s); err != nil {
		t.Fatal(err.Error())
	}
	result := `
// Validate validates the struct
func (m *Message) Validate() error {
	return validator.NewMulti(validator.Date(m.CreatedAt),
		validator.MaxLength(m.Body, 3),
		validator.Maximum(float64(m.Age), 100.000000),
		validator.MinLength(m.Body, 2),
		validator.OneOf(m.StatusConstant, []string{"active", "deleted"}),
		validator.Pattern(m.Body, "^(/[^/]+)+$"))
}`

	code, err := s.GenerateValidators()
	if err != nil {
		t.Fatal(err.Error())
	}

	if result != string(code) {
		t.Fatalf("Schema is not same. Wanted: %s, Get: %s", result, string(code))
	}
}

func TestGenerateFunctions(t *testing.T) {
	var s Schema
	if err := json.Unmarshal([]byte(testJSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	_, err := s.GenerateFunctions()
	if err != nil {
		t.Fatal(err.Error())
	}
}

const testJSON1 = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "Message",
  "description": "Message represents a simple post",
  "type": "object",
  "properties": {
    "Id": {
      "description": "The unique identifier for a message",
      "type": "number",
      "format":"int64"
    },
    "Token": {
      "description": "The token for a message security",
      "type": "string"
    },
    "Body": {
      "description": "The body for a message",
      "type": "string",
      "pattern": "^(/[^/]+)+$",
      "minLength": 2,
      "maxLength": 3
    },
    "Age": {
      "type": "integer",
      "minimum": 0,
      "maximum": 100,
      "exclusiveMaximum": true
    },
    "Enabled": {
      "type": "boolean"
    },
    "StatusConstant": {
      "type": "string",
      "enum": [
        "active",
        "deleted"
      ]
    },
    "CreatedAt": {
      "type": "string",
      "format":"date-time"
    }
  },
  "required": [
    "id",
    "body"
  ]
}
`

// unique identifier of the channel message
// Id int64 `json:"id,string"`
