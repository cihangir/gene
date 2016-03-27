// Package clients generates clients for the generated api
package clients

import (
	"fmt"
	"strings"

	"github.com/cihangir/gene/generators/common"
)

func pathfunc(data *common.TemplateData) string {
	return fmt.Sprintf(
		"%s%s/clients/%s.go",
		data.Settings.Get("fullPathPrefix").(string),
		data.ModuleName,
		strings.ToLower(data.Schema.Title),
	)

}

type Generator struct{}

// Generate generates the client package for given schema
func (c *Generator) Generate(req *common.Req, res *common.Res) error {
	o := &common.Op{
		Name:     "clients",
		Template: ClientsTemplate,
		PathFunc: pathfunc,
		Clear:    true,
	}

	return common.Proces(o, req, res)
}
