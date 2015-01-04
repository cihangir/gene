package errors

var ErrorsTemplate = `
package errs
var (
{{range $defKey, $def := .Definitions}}
    {{range $key, $value := $def.Properties}}
        Err{{$def.Title}}{{$value.Title}}NotSet = errors.New("{{$def.Title}} {{$value.Title}} not set")
    {{end}}
{{end}}
)
`
