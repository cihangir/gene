// Package errors generates the common errors for the modules
package errors

import (
	"fmt"
	"log"
	"strings"

	"github.com/cihangir/gene/generators/common"
	"github.com/kr/pretty"
)

type Generator struct{}

func pathfunc(data *common.TemplateData) string {
	log.Printf("data.Settings.Ge.(string) %# v", pretty.Formatter(data.Settings.Get("fullPathPrefix").(string)))
	return fmt.Sprintf(
		"%s/%s.go",
		data.Settings.Get("fullPathPrefix").(string),
		strings.ToLower(data.Schema.Title),
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
