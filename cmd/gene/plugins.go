package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cihangir/gene/generators/common"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/osext"
)

type Generator struct {
	Plugins map[string]string
	Clients map[string]*plugin.Client
}

func Discover() (*Generator, error) {
	g := &Generator{
		Plugins: make(map[string]string),
		Clients: make(map[string]*plugin.Client),
	}

	if err := g.discover("."); err != nil {
		return nil, err
	}

	exePath, err := osext.Executable()
	if err != nil {
		log.Printf("[ERR] Error loading exe directory: %s", err)
	} else {
		if err := g.discover(filepath.Dir(exePath)); err != nil {
			return nil, err
		}
	}

	for name, path := range g.Plugins {
		g.Clients[name] = g.createPluginClient(path)
	}

	return g, nil
}

func (g *Generator) discover(path string) error {
	var err error

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
	}

	return g.discoverSingle(filepath.Join(path, "gene-*"), &g.Plugins)
}

func (g *Generator) discoverSingle(glob string, m *map[string]string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	if *m == nil {
		*m = make(map[string]string)
	}

	for _, match := range matches {
		file := filepath.Base(match)
		// If the filename has a ".", trim up to there
		if idx := strings.Index(file, "."); idx >= 0 {
			file = file[:idx]
		}

		parts := strings.SplitN(file, "-", 2)
		if len(parts) != 2 {
			continue
		}

		log.Printf("[DEBUG] Discovered plugin: %s = %s", parts[1], match)
		(*m)[parts[1]] = match
	}

	return nil
}

func (g *Generator) createPluginClient(path string) *plugin.Client {
	config := &plugin.ClientConfig{
		Cmd:             pluginCmd(path),
		HandshakeConfig: common.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			// client wont use underlying plugin for serving, so a default empty plugin will work
			"generate": &common.GeneratorPlugin{},
		},
	}

	return plugin.NewClient(config)
}

func pluginCmd(path string) *exec.Cmd {
	cmdPath := ""

	// If the path doesn't contain a separator, look in the same
	// directory as the gene executable first.
	if !strings.ContainsRune(path, os.PathSeparator) {
		exePath, err := osext.Executable()
		if err == nil {
			temp := filepath.Join(
				filepath.Dir(exePath),
				filepath.Base(path))

			if _, err := os.Stat(temp); err == nil {
				cmdPath = temp
			}
		}

		// If we still haven't found the executable, look for it
		// in the PATH.
		if v, err := exec.LookPath(path); err == nil {
			cmdPath = v
		}
	}

	// If we still don't have a path, then just set it to the original
	// given path.
	if cmdPath == "" {
		cmdPath = path
	}

	// Build the command to execute the plugin
	return exec.Command(cmdPath)
}
