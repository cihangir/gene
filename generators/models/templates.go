package models

// PackageTemplate holds the template for the packages of the models
var PackageTemplate = `// Generated struct for {{.}}.
package {{.}}
`

// StructTemplate holds the template for the structs of the models
var StructTemplate = `
{{AsComment .Definition.Description}}
type {{ToUpperFirst .Name}} {{goType .Definition}}
`
