package errors

// ErrorsTemplate holds the template for the errors package
var ErrorsTemplate = `
package errs
var (
{{$moduleName := ToUpperFirst .Title}}
{{range $key, $value := .Properties}}
    Err{{$moduleName}}{{ToUpperFirst $key}}NotSet = errors.New("{{$moduleName}}.{{ToUpperFirst $key}} not set")
{{end}}
)
`
