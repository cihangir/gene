package modules

import (
	"fmt"
	"text/template"

	"bytes"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

// GenerateMainFile handles the main file generation for persistent
// connection rpc server
func (m *Module) GenerateMainFile(rootPath string) error {

	moduleName := stringext.ToLowerFirst(
		m.schema.Title,
	)

	mainFilePath := fmt.Sprintf(
		"%s/%s/main.go",
		rootPath,
		fmt.Sprintf(moduleFolderStucture[0], moduleName),
	)

	f, err := generateMainFile(m.schema)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(mainFilePath, f)
}

func generateMainFile(s *schema.Schema) ([]byte, error) {
	const templateName = "mainfile.tmpl"
	temp := template.New(templateName).Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(MainFileTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, templateName, s); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// MainFileTemplate holds the template for the main file generation
var MainFileTemplate = `
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

	rpcwrap.ServeCustomRPC(
		Mux,
		server,
		false,  // use auth
		"json", // codec name
		jsonrpc.NewServerCodec,
	)

	rpcwrap.ServeHTTPRPC(
		Mux,                    // httpmuxer
		server,                 // rpcserver
		"http_json",            // codec name
		jsonrpc.NewServerCodec, // jsoncodec
		ContextCreator,         // contextCreator
	)

	fmt.Println("Server listening on 3000")
	http.ListenAndServe("localhost:3000", Mux)
}`
