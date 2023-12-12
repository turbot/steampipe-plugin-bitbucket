package bitbucket

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type bitbucketConfig struct {
	Username *string `hcl:"username"`
	Password *string `hcl:"password"`
	BaseUrl  *string `hcl:"base_url"`
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
