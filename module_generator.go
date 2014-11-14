package gene

import "fmt"

type Module struct {
	schema *Schema
}

func NewModule(s *Schema) *Module {
	return &Module{schema: s}
}

func (m *Module) Create() error {
	if err := EnsureFolders(
		"./", // root folder
		createModuleStructure(
			lowFirst(m.schema.Title),
		),
	); err != nil {
		return err
	}

	if err := m.GenerateHandlers(); err != nil {
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
