package kit

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

func GenerateInstrumenting(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	temp := template.New("kit.tmpl").Funcs(context.TemplateFuncs)
	if _, err := temp.Parse(InstrumentingTemplate); err != nil {
		return nil, err
	}

	outputs := make([]common.Output, 0)

	var buf bytes.Buffer
	if err := temp.ExecuteTemplate(&buf, "kit.tmpl", nil); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(
		"%s/kitworker/%s.go",
		context.Config.Target,
		"instrumenting",
	)

	api, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	outputs = append(outputs, common.Output{
		Content: api,
		Path:    path,
	})

	return outputs, nil
}
