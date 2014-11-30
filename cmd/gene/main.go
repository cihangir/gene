package main

import (
	"flag"
	"fmt"
	"log"

	"bitbucket.org/cihangirsavas/gene/generators/modules"
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
	fmt.Println("m-->", m.Create())
}
