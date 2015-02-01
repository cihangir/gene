// Package main provides cli for gene package
package main

import (
	"log"

	"github.com/cihangir/gene/generators/modules"
	"github.com/koding/multiconfig"

	_ "github.com/cihangir/govalidator"
	_ "github.com/cihangir/stringext"
	_ "github.com/koding/logging"
	_ "github.com/lann/squirrel"
	_ "golang.org/x/net/context"
)

func main() {

	conf := &Config{}

	envloader := multiconfig.FlagLoader{}
	if err := envloader.Load(conf); err != nil {
		log.Fatalf("config err:", err.Error())
	}

	m, err := modules.NewFromFile(conf.Schema)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	m.TargetFolderName = conf.Target
	if err := m.Create(); err != nil {
		log.Fatalf(err.Error())
		return
	}

	log.Println("module created with success")
}
