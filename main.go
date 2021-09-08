package main

import (
	"os"

	"github.com/cueblox/blox/plugins"
	"github.com/cueblox/blox/plugins/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Info,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	staticSync := &StaticSync{
		logger: logger,
	}

	dataSync := &DataSync{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	pluginMap := map[string]plugin.Plugin{
		"staticsync": &plugins.PostbuildPlugin{Impl: staticSync},
		"datasync":   &plugins.PostbuildPlugin{Impl: dataSync},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.PostbuildHandshakeConfig,
		Plugins:         pluginMap,
	})
}
