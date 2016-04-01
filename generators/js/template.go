package js

// FunctionsTemplate provides the template for js clients of models
var FunctionsTemplate = `{{$schema := .Schema}}{{$title := $schema.Title}}{{$moduleName := .ModuleName}}module.exports.{{$moduleName}} = {
  // New creates a new local {{ToUpperFirst $title}} js client
  {{ToUpperFirst $title}} = function(){}

  // create validators
  {{ToUpperFirst $title}}.validate = function(data){
    return  null
  }

  // create mapper
  {{ToUpperFirst $title}}.map = function(data){
    return null
  }

{{range $funcKey, $funcValue := $schema.Functions}}
  {{ToUpperFirst $title}}.{{$funcKey}} = function(data, callback) {
  {{GenerateJSValidator $funcValue.Properties.incoming}}
    return processRequest('/{{ToLower $title}}/{{ToLower $funcKey}}', data, callback)
  }
{{end}}
}
`
