// Package modules handle module creation for the given json-schema
package modules

import (
	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/common"

	"github.com/cihangir/schema"
)

type Generator interface {
	Name() string
	Generate(*common.Context, *schema.Schema) ([]common.Output, error)
}

var generators []Generator

func init() {
	generators = []Generator{}
}

// Module holds the required parameters for a module
type Module struct {
	schema *schema.Schema

	context *common.Context

	// TargetFolderName holds the folder name for the module
	TargetFolderName string
}

// NewModule creates a new module with the given Schema
func New(conf *config.Config) (*Module, error) {
	s, err := common.Read(conf.Schema)
	if err != nil {
		return nil, err
	}

	context := common.NewContext()
	context.Config = conf

	return &Module{
		schema:           s.Resolve(s),
		context:          context,
		TargetFolderName: "./",
	}, nil
}

// Create creates the module. While creating the module it handles models,
// handlers, errors, servers, clients and tests generation
func (m *Module) Create() error {
	for _, gen := range generators {
		mgen, err := gen.Generate(m.context, m.schema)
		if err != nil {
			return err
		}

		if err := common.WriteOutput(mgen); err != nil {
			return err
		}
	}

	return nil
}
