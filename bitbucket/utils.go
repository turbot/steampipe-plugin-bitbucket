package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ktrysmt/go-bitbucket"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// // Constants
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

func credentials(d *plugin.QueryData) (string, string) {
	username := os.Getenv("BITBUCKET_USERNAME")
	password := os.Getenv("BITBUCKET_PASSWORD")

	bitbucketConfig := GetConfig(d.Connection)
	if bitbucketConfig.Username != nil {
		username = *bitbucketConfig.Username
	}
	if bitbucketConfig.Password != nil {
		password = *bitbucketConfig.Password
	}
	return username, password
}

// authGet issues an authenticated GET. `client.HttpClient` is a bare *http.Client that carries no
// credentials, so reading through it returns 401 for every private resource; the tables that reach
// the API directly must attach basic auth themselves.
func authGet(ctx context.Context, d *plugin.QueryData, urlStr string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	username, password := credentials(d)
	req.SetBasicAuth(username, password)
	req.Header.Set("Accept", "application/json")

	resp, err := connect(ctx, d).HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%s", resp.Status)
	}
	return resp, nil
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

// page walks a paged Bitbucket collection, calling decode once per page; decode returns the next
// page's URL and whether to continue.
//
// A 404 or 403 yields no rows, matching how the existing tables handle them.
func page(
	ctx context.Context,
	d *plugin.QueryData,
	urlStr string,
	decode func(resp *http.Response) (string, bool, error),
) error {
	return walkPages(ctx, d, urlStr, true, decode)
}

// pageReportingForbidden is page, except that a 403 is returned as an error, for tables where a
// missing scope must not read as an empty collection.
func pageReportingForbidden(
	ctx context.Context,
	d *plugin.QueryData,
	urlStr string,
	decode func(resp *http.Response) (string, bool, error),
) error {
	return walkPages(ctx, d, urlStr, false, decode)
}

func walkPages(
	ctx context.Context,
	d *plugin.QueryData,
	urlStr string,
	forbiddenIsEmpty bool,
	decode func(resp *http.Response) (string, bool, error),
) error {
	for urlStr != "" {
		resp, err := authGet(ctx, d, urlStr)
		if err != nil {
			if isNotFoundError(err) || (forbiddenIsEmpty && isForbiddenError(err)) {
				return nil
			}
			return err
		}
		next, more, err := decode(resp)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		if !more || d.RowsRemaining(ctx) == 0 {
			return nil
		}
		urlStr = next
	}
	return nil
}
