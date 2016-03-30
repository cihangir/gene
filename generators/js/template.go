package js

// FunctionsTemplate provides the template for js clients of models
var FunctionsTemplate = `{{$schema := .Schema}}{{$title := $schema.Title}}module.exports.{{.ModuleName}} = {
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
    // data should be type of {{Argumentize $funcValue.Properties.incoming}}
{{if Equal "array" $funcValue.Properties.incoming.Type}}{{$incoming := index $funcValue.Properties.incoming.Items 0}}
    for (var i = 0; i < data.length; i++){
      if(err = {{$incoming.Title}}.validate(data[i])) {
        return callback(err, null)
      }
    }
{{else}}
    if(err = {{ToUpperFirst $funcValue.Properties.incoming.Title}}.validate(data)) {
      return callback(err, null)
    }
{{end}}

    // send request to the server
    // we got the response
    var res = {}
    // response should be type of {{Argumentize $funcValue.Properties.outgoing}}x{{if Equal "array" $funcValue.Properties.outgoing.Type}}{{$outgoing := index $funcValue.Properties.outgoing.Items 0}}
    res = res.map(function(datum) {
      return {{ToUpperFirst $outgoing.Title}}.map(datum);
    });{{else}}
    res = {{ToUpperFirst $funcValue.Properties.outgoing.Title}}.map(res){{end}}
    callback(null, res)
  }
{{end}}
}
`
