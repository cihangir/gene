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

func (m *Module) GenerateMainFile(rootPath string) error {

	moduleName := stringext.ToLowerFirst(
		m.schema.Title,
	)

	mainFilePath := fmt.Sprintf(
		"%s%s/cmd/%s/main.go",
		rootPath,
		fmt.Sprintf(moduleFolderStucture[0], moduleName),
		moduleName,
	)

	f, err := generateMainFile(m.schema)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(mainFilePath, f)
}

func generateMainFile(s *schema.Schema) ([]byte, error) {
	const templateName = "mainfile.tmpl"
	temp := template.New(templateName)
	temp.Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	})
	_, err := temp.Parse(MainFileTemplate)
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

var MainFileTemplate string = `
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

func main() {
	{{ToLower .Title}} := new({{ToLower .Title}}api.{{.Title}})

	server := rpcplus.NewServer()
	server.Register({{ToLower .Title}})

	mux := http.NewServeMux()

	contextCreator := func(req *http.Request) context.Context {
		return context.Background()
	}

	rpcwrap.ServeHTTPRPC(
		mux,                    // httpmuxer
		server,                 // rpcserver
		"json",                 // codec name
		jsonrpc.NewServerCodec, // jsoncodec
		contextCreator,         // contextCreator
	)

	fmt.Println("Server listening on 3000")
	http.ListenAndServe("localhost:3000", mux)
}`
