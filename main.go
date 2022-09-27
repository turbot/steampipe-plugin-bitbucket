package main

import (
	"github.com/turbot/steampipe-plugin-bitbucket/bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: bitbucket.Plugin})
}
