package freshservice

import (
    "github.com/turbot/steampipe-plugin-sdk/plugin"
    "github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type PluginConfig struct {
    BaseUrl *string `cty:"base_url"`
    Token   *string `cty:"token"`
}

var ConfigSchema = map[string]*schema.Attribute{
    "baseurl": {
        Type: schema.TypeString,
    },
    "token": {
        Type: schema.TypeString,
    },
}

func ConfigInstance() interface{} {
    return &PluginConfig{}
}

func GetConfig(connection *plugin.Connection) PluginConfig {
    if connection == nil || connection.Config == nil {
        return PluginConfig{}
    }

    config, _ := connection.Config.(PluginConfig)
    return config
}
