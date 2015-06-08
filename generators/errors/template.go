package errors

// ErrorsTemplate holds the template for the errors package
var ErrorsTemplate = `
package errs
var (
{{$moduleName := ToUpperFirst .Schema.Title}}
{{range $key, $value := .Schema.Properties}}
    Err{{$moduleName}}{{ToUpperFirst $key}}NotSet = errors.New("{{$moduleName}}.{{ToUpperFirst $key}} not set")
{{end}}
)
`
