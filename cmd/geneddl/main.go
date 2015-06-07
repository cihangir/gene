// Package main provides cli for gene ddl package
package main

import (
	"log"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/generators/sql/definitions"
	"github.com/koding/multiconfig"
)

type Config struct {
	Schema string `required:"true"`
	Target string `default:"ddl"`
	DDL    definitions.Generator
}

func main() {
	conf := &Config{}

	loader := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},  // assign default values
		&multiconfig.FlagLoader{}, // read flag params
		&multiconfig.EnvironmentLoader{},
	)

	if err := loader.Load(conf); err != nil {
		log.Fatalf("config read err: %s", err.Error())
	}

	if err := (&multiconfig.RequiredValidator{}).Validate(conf); err != nil {
		log.Fatalf("validation err: %s", err.Error())
	}

	c := common.NewContext()
	c.Config.Schema = conf.Schema
	c.Config.Target = conf.Target
	c.FieldNameFunc = definitions.GetFieldNameFunc(conf.DDL.FieldNameCase)

	s, err := common.Read(c.Config.Schema)
	if err != nil {
		log.Fatalf("schema read err: %s", err.Error())
	}

	output, err := conf.DDL.Generate(c, s)
	if err != nil {
		log.Fatal("geneddl err: %s", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	log.Println("module created with success")
}
