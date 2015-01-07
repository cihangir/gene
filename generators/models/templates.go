package models

var PackageTemplate = `// Generated struct for {{.}}.
package {{.}}
`
var StructTemplate = `
{{AsComment .Definition.Description}}
type {{ToUpperFirst .Name}} {{goType .Definition}}
`

var FunctionsTemplate = `{{$Name := .Name}}
{{range .Funcs}}
    func ({{Pointerize $Name}} *{{$Name}}){{.}}() {

    }
{{end}}
`
