package plugin

import (
	"net/rpc"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
	"github.com/hashicorp/go-plugin"
)

type Generator interface {
	Generate(context *common.Context, s *schema.Schema) ([]common.Output, error)
}

type Request struct {
	Context *common.Context
	Schema  *schema.Schema
}

type GeneratorRPCServer struct {
	// This is the real implementation
	Impl Generator
}

func (g *GeneratorRPCServer) Generate(args Request, res []common.Output) (err error) {
	res, err = g.Impl.Generate(args.Context, args.Schema)
	return
}

// Here is an implementation that talks over RPC
type GeneratorRPCClient struct{ Client *rpc.Client }

func (g *GeneratorRPCClient) Generate(context *common.Context, s *schema.Schema) (res []common.Output, err error) {
	args := &Request{
		Context: context,
		Schema:  s,
	}

	return res, g.Client.Call("Plugin.Generate", args, &res)
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
