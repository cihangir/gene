// Package main provides models plugin for gene package
package main

import (
	gmodels "github.com/cihangir/gene/generators/models"
	gplugin "github.com/cihangir/gene/plugin"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: gplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"generate": gplugin.NewGeneratorPlugin(&gmodels.Generator{}),
		},
	})
}
