package bitbucket

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type bitbucketConfig struct {
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
	return &bitbucketConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) bitbucketConfig {
	if connection == nil || connection.Config == nil {
		return bitbucketConfig{}
	}
	config, _ := connection.Config.(bitbucketConfig)
	return config
}
