package dockerfiles

import (
	"bytes"
	"fmt"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type Generator struct{}

// Generate generates Dockerfile for given schema
func (c *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	tmpl := template.New("dockerfile.tmpl").Funcs(context.TemplateFuncs)
	if _, err := tmpl.Parse(DockerfileTemplate); err != nil {
		return nil, err
	}

	moduleName := context.ModuleNameFunc(s.Title)
	outputs := make([]common.Output, 0)

	for _, def := range common.SortedObjectSchemas(s.Definitions) {

		var buf bytes.Buffer

		data := struct {
			Target     string
			ModuleName string
			Schema     *schema.Schema
		}{
			Target:     context.Config.Target,
			ModuleName: moduleName,
			Schema:     def,
		}

		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, err
		}

		f, err := format.Source(buf.Bytes())
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf(
			"%sdockerfiles/%s/%s/Dockerfile",
			context.Config.Target,
			moduleName,
			context.FileNameFunc(def.Title),
		)

		outputs = append(outputs, common.Output{Content: f, Path: path, DoNotFormat: true})
	}

	return outputs, nil
}

// DockerfileTemplate holds the template for Dockerfile
var DockerfileTemplate = `
{{$schema := .Schema}}
{{$title := $schema.Title}}

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

RUN go install {{.Target}}workers/{{ToLower $title}}

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/{{ToLower $title}}

# Document that the service listens on port 8080.
EXPOSE 8080
`
