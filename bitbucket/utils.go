package bitbucket

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/ktrysmt/go-bitbucket"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// Constants
const (
	ColumnDescriptionTitle = "Title of the resource."
)

//// HELPER FUNCTIONS

// create service client
func connect(_ context.Context, d *plugin.QueryData) *bitbucket.Client {
	username := os.Getenv("BITBUCKET_USERNAME")
	password := os.Getenv("BITBUCKET_PASSWORD")
	baseurl := os.Getenv("BITBUCKET_API_BASE_URL")

	// Get connection config for plugin
	bitbucketConfig := GetConfig(d.Connection)
	if bitbucketConfig.Username != nil {
		username = *bitbucketConfig.Username
	}
	if bitbucketConfig.Password != nil {
		password = *bitbucketConfig.Password
	}
	if bitbucketConfig.BaseUrl != nil {
		baseurl = *bitbucketConfig.BaseUrl
	}

	if username == "" {
		panic("'username' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if password == "" {
		panic("'password' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	client := bitbucket.NewBasicAuth(username, password)

	// For private bitbucket setup
	if baseurl != "" {
		client.SetApiBaseURL(baseurl)
	}
	return client
}

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

// decode API raw response
func decodeResponse(resp *http.Response, v interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

// decodeJson(apiResponse, responseStruct):: converts raw apiResponse to required output struct
func decodeJson(response interface{}, respObject interface{}) error {
	resp, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, respObject)
	if err != nil {
		return err
	}
	return nil
}

// resource is not found error handling predicate
func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "404")
}

// User don't have required access to all the api on resource
func isForbiddenError(err error) bool {
	return strings.Contains(err.Error(), "403")
}

type ListResponse struct {
	Page     int    `json:"page,omitempty"`
	Pagelen  int    `json:"pagelen,omitempty"`
	MaxDepth int    `json:"maxDepth,omitempty"`
	Size     int    `json:"size,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}
