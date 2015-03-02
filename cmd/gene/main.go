// Package main provides cli for gene package
package main

import (
	"log"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/modules"
	"github.com/koding/multiconfig"

	_ "github.com/cihangir/govalidator"
	_ "github.com/cihangir/stringext"
	_ "github.com/koding/logging"
	_ "github.com/lann/squirrel"
	_ "golang.org/x/net/context"
)

func main() {
	conf := &config.Config{}

	loader := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},  // assign default values
		&multiconfig.FlagLoader{}, // read flag params
	)

	if err := loader.Load(conf); err != nil {
		log.Fatalf("config read err:", err.Error())
	}

	m, err := modules.NewFromFile(conf.Schema)
	if err != nil {
		log.Fatalf("err while reading schema", err.Error())
	}

	m.TargetFolderName = conf.Target
	if err := m.Create(); err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("module created with success")
}
