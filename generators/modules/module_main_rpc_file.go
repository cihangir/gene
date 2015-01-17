package modules

import (
	"fmt"
	"text/template"

	"bytes"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

// GenerateRPCMainFile handles the main file generation for persistent
// connection rpc server
func (m *Module) GenerateRPCMainFile(rootPath string) error {

	moduleName := stringext.ToLowerFirst(
		m.schema.Title,
	)

	mainFilePath := fmt.Sprintf(
		"%s/%s%srpc/main.go",
		rootPath,
		fmt.Sprintf(moduleFolderStucture[0], moduleName),
		moduleName,
	)

	f, err := generateRPCMainFile(m.schema)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(mainFilePath, f)
}

func generateRPCMainFile(s *schema.Schema) ([]byte, error) {
	const templateName = "mainfile.tmpl"
	temp := template.New(templateName).Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(MainRPCFileTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, templateName, s); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// MainRPCFileTemplate holds the template for the main file generation
var MainRPCFileTemplate = `
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

	rpcwrap.ServeCustomRPC(
		mux,
		server,
		false,  // use auth
		"json", // codec name
		jsonrpc.NewServerCodec,
	)

	fmt.Println("Server listening on 3000")
	http.ListenAndServe("localhost:3000", mux)
}`
