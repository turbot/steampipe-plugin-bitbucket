package bitbucket

import (
	"context"
	"os"
	"strings"

	"github.com/ktrysmt/go-bitbucket"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// create service client
func connect(ctx context.Context, d *plugin.QueryData) *bitbucket.Client {
	username := os.Getenv("BITBUCKET_USERNAME")
	password := os.Getenv("BITBUCKET_PASSWORD")

	// Get connection config for plugin
	bitbucketConfig := GetConfig(d.Connection)
	if &bitbucketConfig != nil {
		if bitbucketConfig.Username != nil {
			username = *bitbucketConfig.Username
		}
		if bitbucketConfig.Password != nil {
			password = *bitbucketConfig.Password
		}
	}

	if password == "" {
		panic("'password' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	client := bitbucket.NewBasicAuth(username, password)
	return client
}

//// HELPER FUNCTIONS

func parseRepoFullName(fullName string) (string, string) {
	owner := ""
	repo := ""
	s := strings.Split(fullName, "/")
	owner = s[0]
	if len(s) > 1 {
		repo = s[1]
	}
	return owner, repo
}

func repositoryFullNameQual(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return d.KeyColumnQuals["repository_full_name"].GetStringValue(), nil
}
