package modules

import (
	"fmt"

	"bitbucket.org/cihangirsavas/gene/generators/folders"
	"bitbucket.org/cihangirsavas/gene/generators/handlers"
	"bitbucket.org/cihangirsavas/gene/schema"

	"bitbucket.org/cihangirsavas/gene/stringext"
)

type Module struct {
	schema *schema.Schema
}

func NewModule(s *schema.Schema) *Module {
	return &Module{schema: s}
}

func (m *Module) Create() error {
	rootPath := "./"

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
		createModuleStructure(stringext.ToLowerFirst(
			m.schema.Title,
		)),
	); err != nil {
		return err
	}

	if err := handlers.Generate(rootPath, m.schema.Title); err != nil {
		return err
	}

	if err := m.GenerateMainFile(rootPath); err != nil {
		return err
	}

	return nil
}

var moduleFolderStucture = []string{
	"gene/modules/%[1]s",
	"gene/modules/%[1]s/api",
	"gene/modules/%[1]s/%[1]s",
	"gene/modules/%[1]s/cmd",
	"gene/modules/%[1]s/tests",
	"gene/modules/%[1]s/errors",
	"gene/modules/%[1]s/handlers",
}

func createModuleStructure(name string) []string {
	modified := make([]string, len(moduleFolderStucture))
	for i, str := range moduleFolderStucture {
		modified[i] = fmt.Sprintf(str, name)
	}

	return modified
}
