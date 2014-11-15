package models

var PackageTemplate = `// Generated struct for {{.}}.
package {{.}}
`
var StructTemplate = `
{{AsComment .Definition.Description}}
type {{ToUpperFirst .Name}} {{goType .Definition}}
`

var ValidatorsTemplate = `
// Validate validates the struct
func ({{ToLowerFirst .Name}} *{{.Name}}) Validate() error {
{{GenerateValidator .Definition}}
}
`
