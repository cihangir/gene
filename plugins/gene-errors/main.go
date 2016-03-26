// Package main provides errors plugin for gene package
package main

import (
	gerrors "github.com/cihangir/gene/generators/errors"
	gplugin "github.com/cihangir/gene/plugin"
	"github.com/hashicorp/go-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: gplugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"generate": gplugin.NewGeneratorPlugin(&gerrors.Generator{}),
		},
	})
}
