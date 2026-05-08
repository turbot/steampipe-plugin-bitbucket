package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ktrysmt/go-bitbucket"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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

	client, err := bitbucket.NewBasicAuth(username, password)
	if err != nil {
		panic(fmt.Sprintf("failed to create bitbucket client: %s", err))
	}

	// Wrap the HTTP client so we honor 429 Retry-After from Bitbucket Cloud.
	// The upstream library does not retry on rate limits (ktrysmt/go-bitbucket#295).
	if client.HttpClient == nil {
		client.HttpClient = &http.Client{}
	}
	base := client.HttpClient.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	client.HttpClient.Transport = &retryTransport{base: base, maxRetries: 3, maxWait: 60 * time.Second}

	// For private bitbucket setup
	if baseurl != "" {
		parsed, err := url.Parse(baseurl)
		if err != nil {
			panic(fmt.Sprintf("invalid base_url %q: %s", baseurl, err))
		}
		client.SetApiBaseURL(*parsed)
	}
	return client
}

// retryTransport retries idempotent requests once Bitbucket returns 429
// Too Many Requests, honoring the Retry-After header (seconds or HTTP-date)
// and falling back to exponential backoff when absent. Caps the per-attempt
// wait at maxWait so a long server-suggested delay does not stall a query;
// callers can rerun once the rolling-hour window clears.
type retryTransport struct {
	base       http.RoundTripper
	maxRetries int
	maxWait    time.Duration
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only retry methods that are safe to replay. Body-bearing requests
	// would need GetBody to rewind; the plugin only issues GETs today, but
	// guard explicitly so that any future POST does not silently break.
	canRetry := req.Method == http.MethodGet || req.Method == http.MethodHead

	var resp *http.Response
	var err error
	for attempt := 0; ; attempt++ {
		resp, err = t.base.RoundTrip(req)
		if err != nil || resp.StatusCode != http.StatusTooManyRequests || !canRetry || attempt >= t.maxRetries {
			return resp, err
		}

		wait := parseRetryAfter(resp.Header.Get("Retry-After"))
		if wait <= 0 {
			wait = time.Duration(1<<attempt) * time.Second
		}
		if wait > t.maxWait {
			// Surface the 429 so the query fails fast rather than blocking
			// the connection for an unbounded period.
			return resp, nil
		}

		// Drain and close the response so the underlying connection can be reused.
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		select {
		case <-time.After(wait):
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}
}

func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
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
