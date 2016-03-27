package common

import "github.com/cihangir/gene/config"

type Context struct {
	Config *config.Config
}

func NewContext() *Context {
	return &Context{
		Config: &config.Config{
			Target: "./",
			Generators: []string{
				"ddl", "rows", "kit", "errors",
				"dockerfiles", "clients", "tests",
				"functions", "models",
			},
		},
	}
}
