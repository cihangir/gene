package validators

import (
	"fmt"
	"sort"
	"strings"

	"bitbucket.org/cihangirsavas/gene/schema"
	"bitbucket.org/cihangirsavas/gene/stringext"
)

func GenerateValidator(p *schema.Schema) string {
	validators := make([]string, 0)
	// schemaName := p.Title
	schemaFirstChar := stringext.Pointerize(p.Title)

	for key, property := range p.Properties {
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

			switch property.Format {
			case "date-time":
				// _, err := time.Parse(time.RFC3339, s)
				validator := fmt.Sprintf("validator.Date(%s.%s)", schemaFirstChar, key)
				validators = append(validators, validator)
			}

		case "integer":

			// todo implement exclusive min/max

			if property.Minimum != 0 {
				validator := fmt.Sprintf("validator.Minimum(float64(%s.%s), %f)", schemaFirstChar, key, property.Minimum)
				validators = append(validators, validator)
			}

			if property.Maximum != 0 {
				validator := fmt.Sprintf("validator.Maximum(float64(%s.%s), %f)", schemaFirstChar, key, property.Maximum)
				validators = append(validators, validator)
			}

			// multipleOf:
			if property.MultipleOf != 0 {
				validator := fmt.Sprintf("validator.MultipleOf(float64(%s.%s), %f)", schemaFirstChar, key, property.MultipleOf)
				validators = append(validators, validator)
			}
		}
	}

	// keep the order of validators in same with every call
	sslice := sort.StringSlice(validators)
	sslice.Sort()
	validatorsStr := strings.Join(sslice, ",\n")
	a := fmt.Sprintf("return validator.NewMulti(%s)", validatorsStr)
	return a
}
