package common

import "github.com/cihangir/gene/config"

type Context struct {
	Config *config.Config
}

func NewContext() *Context {
	return &Context{
		Config: &config.Config{},
	}
}
