package plugin

import (
	"net/rpc"

	"github.com/cihangir/gene/generators/common"
	"github.com/hashicorp/go-plugin"
)

type GeneratorRPCServer struct{ Impl common.Generator }

func (g *GeneratorRPCServer) Generate(req *common.Req, res *common.Res) error {
	return g.Impl.Generate(req, res)
}

// Here is an implementation that talks over RPC
type GeneratorRPCClient struct{ Client *rpc.Client }

func (g *GeneratorRPCClient) Generate(req *common.Req, res *common.Res) error {
	return g.Client.Call("Plugin.Generate", req, res)
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GENE_PLUGIN",
	MagicCookieValue: "gene-cookie",
}

type GeneratorPlugin struct{ generator common.Generator }

func NewGeneratorPlugin(g common.Generator) *GeneratorPlugin {
	return &GeneratorPlugin{generator: g}
}

func (g *GeneratorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GeneratorRPCServer{Impl: g.generator}, nil
}

func (g *GeneratorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GeneratorRPCClient{Client: c}, nil
}
