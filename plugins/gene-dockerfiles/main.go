// Package main provides dockerfiles plugin for gene package
package main

import (
	gdockerfiles "github.com/cihangir/gene/generators/dockerfiles"
	gplugin "github.com/cihangir/gene/plugin"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: gplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"generate": gplugin.NewGeneratorPlugin(&gdockerfiles.Generator{}),
		},
	})
}
