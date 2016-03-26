// Package main provides rows plugin for gene package
package main

import (
	gplugin "github.com/cihangir/gene/plugin"
	grows "github.com/cihangir/generows"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: gplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"generate": gplugin.NewGeneratorPlugin(&grows.Generator{}),
		},
	})
}
