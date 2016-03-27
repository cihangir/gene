package common

import "github.com/cihangir/schema"

type Req struct {
	Schema    *schema.Schema
	SchemaStr string
	Context   *Context
}

type Res struct {
	Output []Output
}

type Generator interface {
	Generate(*Req, *Res) error
}
