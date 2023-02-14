package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"steampipe-plugin-freshservice/freshservice"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: freshservice.Plugin})
}
