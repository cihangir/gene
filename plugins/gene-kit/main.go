// Package main provides kit plugin for gene package
package main

import (
	gkit "github.com/cihangir/gene/generators/kit"
	gplugin "github.com/cihangir/gene/plugin"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: gplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"generate": gplugin.NewGeneratorPlugin(&gkit.Generator{}),
		},
	})
}
