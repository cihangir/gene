// Package modules handle module creation for the given json-schema
package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/clients"
	"github.com/cihangir/gene/generators/common"
	gerr "github.com/cihangir/gene/generators/errors"
	"github.com/cihangir/gene/generators/folders"
	"github.com/cihangir/gene/generators/functions"
	"github.com/cihangir/gene/generators/mainfile"
	"github.com/cihangir/gene/generators/models"
	"github.com/cihangir/gene/generators/scanners/rows"
	"github.com/cihangir/gene/generators/sql/definitions"
	"github.com/cihangir/gene/generators/sql/statements"
	"github.com/cihangir/gene/generators/tests"
	"github.com/cihangir/gene/writers"

	"github.com/cihangir/gene/helpers"
	"github.com/cihangir/schema"
	"gopkg.in/yaml.v2"
)

type Generator interface {
	Name() string
	Generate(*config.Context, *schema.Schema) ([]common.Output, error)
}

var generators []Generator

func init() {
	generators = []Generator{
		statements.New(),
		models.New(),
		rows.New(),
		gerr.New(),
		mainfile.New(),
		clients.New(),
		tests.New(),
		functions.New(),
		definitions.New(),
		// js.New(),
	}
}

// Module holds the required parameters for a module
type Module struct {
	schema *schema.Schema

	context *config.Context

	// TargetFolderName holds the folder name for the module
	TargetFolderName string
}

// NewModule creates a new module with the given Schema
func New(conf *config.Config) (*Module, error) {
	s, err := read(conf)
	if err != nil {
		return nil, err
	}

	context := config.NewContext()
	context.Config = conf

	return &Module{
		schema:           s.Resolve(s),
		context:          context,
		TargetFolderName: "./",
	}, nil
}

// NewFromFile reads the given file and creates a new module out of it
func read(config *config.Config) (*schema.Schema, error) {
	fileContent, err := helpers.ReadFile(config.Schema)
	if err != nil {
		return nil, err
	}

	return unmarshall(config.Schema, fileContent)
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

	// mgenerator, err := models.New(m.context, m.schema)
	// if err != nil {
	// 	return err
	// }
	//
	for _, gen := range generators {
		mgen, err := gen.Generate(m.context, m.schema)
		if err != nil {
			return err
		}

		for _, file := range mgen {
			// do not write empty files
			if len(file.Content) == 0 {
				continue
			}

			if file.DoNotFormat {
				if err := writers.Write(file.Path, file.Content); err != nil {
					return err
				}
			} else {
				if err := writers.WriteFormattedFile(file.Path, file.Content); err != nil {
					return err
				}
			}

		}
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
