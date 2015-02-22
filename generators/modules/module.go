// Package modules handle module creation for the given json-schema
package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cihangir/gene/generators/clients"
	gerr "github.com/cihangir/gene/generators/errors"
	"github.com/cihangir/gene/generators/folders"
	"github.com/cihangir/gene/generators/functions"
	"github.com/cihangir/gene/generators/js"
	"github.com/cihangir/gene/generators/models"
	"github.com/cihangir/gene/generators/tests"

	"github.com/cihangir/gene/helpers"
	"github.com/cihangir/schema"
	"gopkg.in/yaml.v2"
)

// Module holds the required parameters for a module
type Module struct {
	schema *schema.Schema

	// TargetFolderName holds the folder name for the module
	TargetFolderName string
}

// NewModule creates a new module with the given Schema
func NewModule(s *schema.Schema) *Module {
	return &Module{
		schema:           s.Resolve(s),
		TargetFolderName: "./",
	}
}

// NewFromFile reads the given file and creates a new module out of it
func NewFromFile(path string) (*Module, error) {
	fileContent, err := helpers.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s, err := unmarshall(path, fileContent)
	if err != nil {
		return nil, err
	}

	return NewModule(s), nil
}

func unmarshall(path string, fileContent []byte) (*schema.Schema, error) {
	s := &schema.Schema{}

	// Choose what while is passed
	switch filepath.Ext(path) {
	case ".toml":
		if err := toml.Unmarshal(fileContent, s); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(fileContent, s); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(fileContent, s); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unmarshal not implemented")
	}

	return s, nil
}

// Create creates the module. While creating the module it handles models,
// handlers, errors, servers, clients and tests generation
func (m *Module) Create() error {
	rootPath := m.TargetFolderName

	// first ensure that we have the correct folder structure for our system
	if err := folders.EnsureFolders(
		rootPath, // root folder
		folders.FolderStucture,
	); err != nil {
		return err
	}

	// create the module folder structure
	if err := folders.EnsureFolders(
		rootPath, // root folder
		createModuleStructure(strings.ToLower(
			m.schema.Title,
		)),
	); err != nil {
		return err
	}

	if err := models.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := models.GenerateStatements(rootPath, m.schema); err != nil {
		return err
	}

	if err := gerr.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := m.GenerateMainFile(rootPath); err != nil {
		return err
	}

	if err := clients.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := tests.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := functions.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := js.Generate(rootPath, m.schema); err != nil {
		return err
	}

	return nil
}

var moduleFolderStucture = []string{
	"cmd/%[1]s/",
	"workers/%[1]s",
	"workers/%[1]s/api",
	"workers/%[1]s/tests",
	"workers/%[1]s/js",
	"workers/%[1]s/errors",
	"workers/%[1]s/clients",
}

func createModuleStructure(name string) []string {
	modified := make([]string, len(moduleFolderStucture))
	for i, str := range moduleFolderStucture {
		modified[i] = fmt.Sprintf(str, name)
	}

	return modified
}
