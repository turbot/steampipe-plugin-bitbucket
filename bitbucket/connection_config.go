package bitbucket

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type githubConfig struct {
	Username *string `cty:"username"`
	Password *string `cty:"password"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"username": {
		Type: schema.TypeString,
	},
	"password": {
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
