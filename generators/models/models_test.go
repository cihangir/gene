package models

import (
	"encoding/json"
	"strings"

	"testing"

	"github.com/cihangir/gene/generators/folders"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/gene/writers"
)

func TestGenerateModel(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	model, err := GenerateModel(&s)
	if err != nil {
		t.Fatal(err.Error())
	}

	folders.EnsureFolders("/tmp/", folders.FolderStucture)
	fileName := "/tmp/gene/models/" + s.Title + ".go"

	err = writers.WriteFormattedFile(fileName, model)
	if err != nil {
		t.Fatal(err.Error())
	}

}

func TestGenerateSchema(t *testing.T) {
	var s schema.Schema
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

	code, err := GenerateSchema(&s)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result != string(code) {
		// fmt.Printf("foo %# v", pretty.Formatter(difflib.Diff([]string{result}, []string{string(code)})))
		t.Fatalf("Schema is not same. Wanted: %s, Get: %s", result, string(code))
	}
}

func TestGenerateValidators(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
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

	code, err := GenerateValidators(&s)
	if err != nil {
		t.Fatal(err.Error())
	}

	if result != string(code) {
		t.Fatalf("Schema is not same. Wanted: %s, Get: %s", result, string(code))
	}
}

func TestGenerateFunctions(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	_, err := GenerateFunctions(&s)
	if err != nil {
		t.Fatal(err.Error())
	}
}
