// Package main provides cli for gene package
package main

import (
	"log"

	"github.com/cihangir/gene/generators/clients"
	"github.com/cihangir/gene/generators/common"
	gerr "github.com/cihangir/gene/generators/errors"
	"github.com/cihangir/gene/generators/functions"
	"github.com/cihangir/gene/generators/mainfile"
	"github.com/cihangir/gene/generators/models"
	"github.com/cihangir/gene/generators/models/rowsscanner"
	"github.com/cihangir/gene/generators/sql/definitions"
	"github.com/cihangir/gene/generators/sql/statements"
	"github.com/koding/multiconfig"

	_ "github.com/cihangir/govalidator"
	_ "github.com/cihangir/stringext"
	_ "github.com/koding/logging"
	_ "github.com/lann/squirrel"
	_ "golang.org/x/net/context"
)

type Config struct {
	// Schema holds the given schema file
	Schema string `required:"true"`

	// Target holds the target folder
	Target string `required:"true" default:"./"`

	DDL    definitions.Generator
	Models models.Generator

	Rows       rows.Generator
	Statements statements.Generator
	Errors     gerr.Generator
	Mainfile   mainfile.Generator
	Clients    clients.Generator
	Functions  functions.Generator
}

func main() {
	conf := &Config{}

	loader := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},  // assign default values
		&multiconfig.FlagLoader{}, // read flag params
	)

	if err := loader.Load(conf); err != nil {
		log.Fatalf("config read err:", err.Error())
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

	s.Resolve(s)

	//
	// generate sql definitions
	//
	c.Config.Target = conf.Target + "db" + "/"
	output, err := conf.DDL.Generate(c, s)
	if err != nil {
		log.Fatal("geneddl err: %s", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate models
	//
	c.Config.Target = conf.Target + "models" + "/"
	output, err = conf.Models.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating models", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate rowsscanner
	//
	c.Config.Target = conf.Target + "models" + "/"
	output, err = conf.Rows.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating rows", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate crud statements
	//
	c.Config.Target = conf.Target + "models" + "/"
	output, err = conf.Statements.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating crud statements", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate errors
	//
	c.Config.Target = conf.Target + "errors" + "/"
	output, err = conf.Errors.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating errors", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate main file
	//
	c.Config.Target = conf.Target + "cmd" + "/"
	output, err = conf.Mainfile.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating main file", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate clients
	//
	c.Config.Target = conf.Target + "workers" + "/"
	output, err = conf.Clients.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating clients", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	//
	// generate exported functions
	//
	c.Config.Target = conf.Target + "workers" + "/"
	output, err = conf.Functions.Generate(c, s)
	if err != nil {
		log.Fatalf("err while generating clients", err.Error())
	}

	if err := common.WriteOutput(output); err != nil {
		log.Fatal("output write err: %s", err.Error())
	}

	log.Println("module created with success")
}
