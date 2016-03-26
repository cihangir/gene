// Package errors generates the common errors for the modules
package errors

import (
	"fmt"
	"strings"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type Generator struct{}

func pathfunc(context *common.Context, def *schema.Schema) string {
	return fmt.Sprintf(
		"%s/%s.go",
		context.Config.Target,
		strings.ToLower(def.Title),
	)
}

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(req *common.Req, res *common.Res) error {
	o := &common.Op{
		Name:     "errors",
		Template: ErrorsTemplate,
		PathFunc: pathfunc,
		Clear:    true,
	}

	return common.Proces(o, req, res)
}
