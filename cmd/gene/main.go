// Package main provides cli for gene package
package main

import (
	"log"

	"github.com/cihangir/gene/generators/clients"
	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/generators/dockerfiles"
	gerr "github.com/cihangir/gene/generators/errors"
	"github.com/cihangir/gene/generators/functions"
	"github.com/cihangir/gene/generators/kit"
	"github.com/cihangir/gene/generators/mainfile"
	"github.com/cihangir/gene/generators/models"
	"github.com/cihangir/gene/generators/sql/statements"
	"github.com/cihangir/geneddl"
	"github.com/cihangir/generows"
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

	// Generators holds the generator names for processing
	Generators []string `default:"ddl,rows,kit,errors,dockerfiles,clients,tests,functions,models"`

	DDL    geneddl.Generator
	Models models.Generator

	Rows        generows.Generator
	Statements  statements.Generator
	Errors      gerr.Generator
	Mainfile    mainfile.Generator
	Clients     clients.Generator
	Functions   functions.Generator
	Dockerfiles dockerfiles.Generator
	// Js         js.Generator
	// Server     server.Generator
	Kit kit.Generator
}

func main() {
	conf := &Config{}

	g, err := Discover()
	if err != nil {
		log.Fatalf("err %# s", err)
	}

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
	c.Config.Target = conf.Target
	c.Config.Generators = conf.Generators

	str, err := common.ReadJSON(conf.Schema)
	if err != nil {
		log.Fatalf("schema read err: %s", err.Error())
	}

	for name, client := range g.Clients {
		log.Print("generating for ", name)

		rpcClient, err := client.Client()
		if err != nil {
			log.Fatalf("couldnt start client", err)
		}
		defer rpcClient.Close()

		raw, err := rpcClient.Dispense("generate")
		if err != nil {
			log.Fatalf("couldnt get the client", err)
		}

		gene := (raw).(common.Generator)
		req := &common.Req{
			SchemaStr: str,
			Context:   c,
		}

		res := &common.Res{}
		err = gene.Generate(req, res)
		if err != nil {
			log.Fatalf("err while generating content for %s, err: %# v", name, err)
		}

		if err := common.WriteOutput(res.Output); err != nil {
			log.Fatal("output write err: %s", err.Error())
		}
	}

	//
	// generate crud statements
	//
	// c.Config.Target = conf.Target + "models" + "/"
	// output, err = conf.Statements.Generate(c, s)
	// if err != nil {
	// 	log.Fatalf("err while generating crud statements", err.Error())
	// }

	// if err := common.WriteOutput(output); err != nil {
	// 	log.Fatal("output write err: %s", err.Error())
	// }

	//
	// generate main file
	//
	// c.Config.Target = conf.Target + "cmd" + "/"
	// output, err = conf.Mainfile.Generate(c, s)
	// if err != nil {
	// 	log.Fatalf("err while generating main file", err.Error())
	// }

	// if err := common.WriteOutput(output); err != nil {
	// 	log.Fatal("output write err: %s", err.Error())
	// }

	//
	// generate clients
	//
	// c.Config.Target = conf.Target + "workers" + "/"
	// output, err = conf.Clients.Generate(c, s)
	// if err != nil {
	// 	log.Fatalf("err while generating clients", err.Error())
	// }

	// if err := common.WriteOutput(output); err != nil {
	// 	log.Fatal("output write err: %s", err.Error())
	// }

	//
	// generate exported functions
	//
	// c.Config.Target = conf.Target + "workers" + "/"
	// output, err = conf.Functions.Generate(c, s)
	// if err != nil {
	// 	log.Fatalf("err while generating clients", err.Error())
	// }

	// if err := common.WriteOutput(output); err != nil {
	// 	log.Fatal("output write err: %s", err.Error())
	// }

	//
	// generate js client functions
	//
	// c.Config.Target = conf.Target + "js" + "/"
	// output, err = conf.Js.Generate(c, s)
	// if err != nil {
	// 	log.Fatalf("err while generating js clients", err.Error())
	// }

	// if err := common.WriteOutput(output); err != nil {
	// 	log.Fatal("output write err: %s", err.Error())
	// }

	//
	// generate api server handlers
	//
	// c.Config.Target = conf.Target + "api" + "/"
	// output, err = conf.Server.Generate(c, s)
	// if err != nil {
	// 	log.Fatalf("err while generating api server", err.Error())
	// }

	// if err := common.WriteOutput(output); err != nil {
	// 	log.Fatal("api output write err: %s", err.Error())
	// }

	//
	// generate kit server handlers
	//

	log.Println("module created with success")
}
