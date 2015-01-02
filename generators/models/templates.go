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
func ({{Pointerize .Name}} *{{.Name}}) Validate() error {
{{GenerateValidator .Definition}}
}
`

var FunctionsTemplate = `{{$Name := .Name}}
{{range .Funcs}}
    func ({{Pointerize $Name}} *{{$Name}}){{.}}() {

    }
{{end}}
`
