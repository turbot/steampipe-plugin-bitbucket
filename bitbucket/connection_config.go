package bitbucket

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type githubConfig struct {
	Username *string `cty:"username"`
	Password *string `cty:"password"`
	BaseUrl  *string `cty:"base_url"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"username": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"base_url": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &githubConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) githubConfig {
	if connection == nil || connection.Config == nil {
		return githubConfig{}
	}
	config, _ := connection.Config.(githubConfig)
	return config
}
