package errors

var ErrorsTemplate = `
package errs
var (
{{range $key, $value := .Properties}}
    Err{{$key}}NotSet = errors.New("{{$key}} not set")
{{end}}
)
`
