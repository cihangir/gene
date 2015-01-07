package validators

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(s *schema.Schema) ([]byte, error) {
	validators := make([]string, 0)
	// schemaName := p.Title
	schemaFirstChar := stringext.Pointerize(s.Title)

	for key, property := range s.Properties {
		switch property.Type {
		case "string":
			if property.MinLength != 0 {
				validator := fmt.Sprintf("validator.MinLength(%s.%s, %d)", schemaFirstChar, key, property.MinLength)
				validators = append(validators, validator)
			}

			if property.MaxLength != 0 {
				validator := fmt.Sprintf("validator.MaxLength(%s.%s, %d)", schemaFirstChar, key, property.MaxLength)
				validators = append(validators, validator)
			}

			if property.Pattern != "" {
				validator := fmt.Sprintf("validator.Pattern(%s.%s, \"%s\")", schemaFirstChar, key, property.Pattern)
				validators = append(validators, validator)
			}

			if len(property.Enum) > 0 {
				validator := fmt.Sprintf("validator.OneOf(%s.%s, []string{\"%s\"})", schemaFirstChar, key, strings.Join(property.Enum, "\",\""))
				validators = append(validators, validator)
			}

			// TODO impplement this one
			switch property.Format {
			case "date-time":
				// _, err := time.Parse(time.RFC3339, s)
				validator := fmt.Sprintf("validator.Date(%s.%s)", schemaFirstChar, key)
				validators = append(validators, validator)
			}

		case "integer", "number":

			// todo implement exclusive min/max

			if property.Minimum != 0 {
				validator := fmt.Sprintf("validator.Min(float64(%s.%s), %f)", schemaFirstChar, key, property.Minimum)
				validators = append(validators, validator)
			}

			if property.Maximum != 0 {
				validator := fmt.Sprintf("validator.Max(float64(%s.%s), %f)", schemaFirstChar, key, property.Maximum)
				validators = append(validators, validator)
			}

			// multipleOf:
			if property.MultipleOf != 0 {
				validator := fmt.Sprintf("validator.MultipleOf(float64(%s.%s), %f)", schemaFirstChar, key, property.MultipleOf)
				validators = append(validators, validator)
			}
		}
	}

	if len(validators) == 0 {
		return nil, nil
	}

	templ := `
// Validate validates the struct
func (%s *%s) Validate() error {
	return validator.NewMulti(%s)
}`

	sslice := sort.StringSlice(validators)
	sslice.Sort()

	res := fmt.Sprintf(
		templ,
		stringext.Pointerize(s.Title),
		s.Title,
		strings.Join(sslice, ",\n"),
	)

	// fmt.Println("res-->", res)
	return writers.Clear(*bytes.NewBufferString(res))
}
