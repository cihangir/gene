package main

import (
	"flag"
	"log"

	"github.com/cihangir/gene/generators/modules"
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
