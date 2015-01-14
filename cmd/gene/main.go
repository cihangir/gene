// Package main provides cli for gene package
package main

import (
	"flag"
	"log"

	"github.com/cihangir/gene/generators/modules"

	_ "github.com/koding/logging"
	_ "github.com/koding/multiconfig"
	_ "github.com/lann/squirrel"
	_ "golang.org/x/net/context"
)

var (
	flagSchemaFile = flag.String("schema", "", "schema content file")
	flagFolder     = flag.String("target", "./", "target directory name")
)

func main() {
	flag.Parse()

	m, err := modules.NewFromFile(*flagSchemaFile)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	m.TargetFolderName = *flagFolder
	if err := m.Create(); err != nil {
		log.Fatalf(err.Error())
		return
	}

	log.Println("module created with success")
}
