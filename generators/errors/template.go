package errors

var ErrorsTemplate = `
package errs
var (
{{range $defKey, $def := .Definitions}}
    {{range $key, $value := $def.Properties}}
        Err{{$def.Title}}{{$key}}NotSet = errors.New("{{$def.Title}}.{{$key}} not set")
    {{end}}
{{end}}
)
`
