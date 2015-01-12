package modules

import (
	"fmt"
	"strings"
	"text/template"

	"bytes"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func (m *Module) GenerateHTTPMainFile(rootPath string) error {

	moduleName := stringext.ToLowerFirst(
		m.schema.Title,
	)

	mainFilePath := fmt.Sprintf(
		"%s%s/cmd/%shttp/main.go",
		rootPath,
		fmt.Sprintf(moduleFolderStucture[0], moduleName),
		moduleName,
	)

	f, err := generateHTTPMainFile(m.schema)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(mainFilePath, f)
}

func generateHTTPMainFile(s *schema.Schema) ([]byte, error) {
	const templateName = "mainfile.tmpl"
	temp := template.New(templateName)
	temp.Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	})
	_, err := temp.Parse(MainHTTPFileTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, templateName, s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var MainHTTPFileTemplate string = `
package main

import (
    "fmt"
    "net/http"
    "github.com/youtube/vitess/go/rpcplus"
    "github.com/youtube/vitess/go/rpcplus/jsonrpc"
    "github.com/youtube/vitess/go/rpcwrap"
)

var (
    Name    = "{{.Title}}"
    VERSION string
)

var ContextCreator = func(req *http.Request) context.Context {
    return context.Background()
}

var Mux = http.NewServeMux()

func main() {

    {{$Name := .Title}}
    server := rpcplus.NewServer()
    {{range $key, $value := .Definitions}}
    server.Register(new({{ToLower $Name}}api.{{$key}}))
    {{end}}

    mux := http.NewServeMux()

    rpcwrap.ServeHTTPRPC(
        Mux,                    // httpmuxer
        server,                 // rpcserver
        "json",                 // codec name
        jsonrpc.NewServerCodec, // jsoncodec
        ContextCreator,         // contextCreator
    )

    fmt.Println("Server listening on 3000")
    http.ListenAndServe("localhost:3000", mux)
}`
