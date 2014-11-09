package gene

import "fmt"

type Module struct {
	schema *Schema
}

func NewModule(s *Schema) *Module {
	return &Module{schema: s}
}

func (m *Module) Create() error {
	return EnsureFolders("./", createModuleStructure(lowFirst(m.schema.Title)))
}

var moduleFolderStucture = []string{
	"gene/modules/%[1]s",
	"gene/modules/%[1]s/%[1]sapi",
	"gene/modules/%[1]s/%[1]s",
	"gene/modules/%[1]s/cmd",
	"gene/modules/%[1]s/%[1]stests",
	"gene/modules/%[1]s/%[1]serrors",
}

func createModuleStructure(name string) []string {
	modified := make([]string, len(moduleFolderStucture))
	for i, str := range moduleFolderStucture {
		modified[i] = fmt.Sprintf(str, name)
	}

	return modified
}
