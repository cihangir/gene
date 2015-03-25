package config

import (
	"strings"
	"text/template"

	"github.com/cihangir/gene/generators/common"
)

// Config holds the config parameters for gene package
type Config struct {
	// Schema holds the given schema file
	Schema string `required:"true"`

	// Target holds the target folder
	Target string `required:"true" default:"./"`

	// Generators holds the generator names for processing
	Generators []string `default:"model,statements,errors,clients,tests,functions"`
}

type Context struct {
	Config *Config

	// Funcs
	ModuleNameFunc func(string) string
	FileNameFunc   func(string) string

	// TemplateFuncs
	TemplateFuncs template.FuncMap
}

func NewContext() *Context {
	return &Context{
		// Funcs
		ModuleNameFunc: strings.ToLower,
		FileNameFunc:   strings.ToLower,

		// TemplateFuncs
		TemplateFuncs: common.TemplateFuncs,
	}
}
