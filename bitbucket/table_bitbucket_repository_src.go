package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

// https://developer.atlassian.com/cloud/bitbucket/rest/api-group-source/#api-repositories-workspace-repo-slug-src-commit-path-get
func tableBitbucketRepositorySrc(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_repository_src",
		Description: "Files and directories in a Bitbucket repository at a commit, and the content of a single file.",
		List: &plugin.ListConfig{
			// Exact cache matching: these quals change what is fetched, so a cached listing must never
			// answer a single-file read or a different version.
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "commit", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "scope_path", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "path", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "max_depth", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
				{Name: "q", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Hydrate: tableBitbucketRepositorySrcList,
		},
		Columns: []*plugin.Column{
			{
				Name:        "repository_full_name",
				Description: "Full name of the repository, e.g. \"my-workspace/my-repo\".",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_full_name"),
			},
			{
				Name:        "commit",
				Description: "The commit hash, branch or tag read, as given. Defaults to HEAD, the head of the main branch.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("commit"),
			},
			{
				Name:        "scope_path",
				Description: "The directory whose contents are listed. Defaults to the repository root.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("scope_path"),
			},
			{
				Name:        "max_depth",
				Description: "How many levels of nested directories to list. Bitbucket defaults to 1.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromQual("max_depth"),
			},
			{
				Name:        "q",
				Description: "Bitbucket filter applied to the listing, e.g. type=\"commit_file\" AND path~\"package.json\".",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("q"),
			},
			{
				Name:        "path",
				Description: "Path of the file or directory from the repository root. Given as a qualifier, reads that one file and its content.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the entry: commit_file or commit_directory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Size of the file in bytes. Null for a directory.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "commit_hash",
				Description: "Hash of the commit the entry was read at. For a single-file read, set only when commit is a full hash.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Commit.Hash").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "mimetype",
				Description: "MIME type Bitbucket infers from the file name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attributes",
				Description: "File attributes such as executable, link, lfs or binary.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "content",
				Description: "Content of the file. Set only when path is given; null for a directory and for binary or non-UTF-8 files.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Path"),
			},
		},
	}
}

func tableBitbucketRepositorySrcList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketRepositorySrcList")

	workspace, repoSlug := parseRepoFullName(d.EqualsQuals["repository_full_name"].GetStringValue())
	if workspace == "" || repoSlug == "" {
		return nil, nil
	}
	req := srcRequest{
		Workspace: workspace,
		RepoSlug:  repoSlug,
		Commit:    d.EqualsQualString("commit"),
		ScopePath: d.EqualsQualString("scope_path"),
		Path:      d.EqualsQualString("path"),
		Query:     d.EqualsQualString("q"),
	}
	if d.EqualsQuals["max_depth"] != nil {
		req.MaxDepth = d.EqualsQuals["max_depth"].GetInt64Value()
	}
	if limit := d.QueryContext.Limit; limit != nil && *limit > 0 && *limit < srcPageLen {
		req.PageLen = *limit
	}

	get := func(urlStr string) (*http.Response, error) { return authGet(ctx, d, urlStr) }
	err := readSrc(get, connect(ctx, d).GetApiBaseURL(), req, func(entry SrcEntry) bool {
		d.StreamListItem(ctx, entry)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		return nil, fmt.Errorf("bitbucket_repository_src: reading %s/%s: %w", workspace, repoSlug, err)
	}
	return nil, nil
}

// Bitbucket's maximum page size for the src listing; larger pages spend fewer of the 1,000 requests
// an hour a token gets.
const srcPageLen = 100

var fullCommitHash = regexp.MustCompile(`^[0-9a-f]{40}$`)

type srcRequest struct {
	Workspace string
	RepoSlug  string
	Commit    string
	ScopePath string
	Path      string
	MaxDepth  int64
	Query     string
	PageLen   int64
}

type SrcList struct {
	ListResponse
	Values []SrcEntry `json:"values"`
}

type SrcEntry struct {
	Path       string    `json:"path"`
	Type       string    `json:"type"`
	Size       *int64    `json:"size"`
	Mimetype   *string   `json:"mimetype"`
	Attributes []string  `json:"attributes"`
	Commit     SrcCommit `json:"commit"`
	Content    *string   `json:"-"`
}

type SrcCommit struct {
	Hash string `json:"hash"`
}

// readSrc lists a directory page by page, or reads one file, calling emit per row until it returns
// false. Every case costs one request per page: Bitbucket answers a directory with a JSON listing
// and a file with its raw bytes, so the response itself says which one the path was.
func readSrc(get func(string) (*http.Response, error), apiBase string, req srcRequest, emit func(SrcEntry) bool) error {
	if req.Path != "" && req.ScopePath != "" {
		return errors.New("give either path (one file) or scope_path (a directory listing), not both")
	}
	commit := req.Commit
	if commit == "" {
		commit = "HEAD"
	}

	if req.Path != "" {
		return readSrcFile(get, srcURL(apiBase, req.Workspace, req.RepoSlug, commit, req.Path, false), req.Path, commit, emit)
	}

	pageLen := req.PageLen
	if pageLen <= 0 {
		pageLen = srcPageLen
	}
	params := url.Values{}
	params.Set("pagelen", strconv.FormatInt(pageLen, 10))
	if req.MaxDepth > 0 {
		params.Set("max_depth", strconv.FormatInt(req.MaxDepth, 10))
	}
	if req.Query != "" {
		params.Set("q", req.Query)
	}
	urlStr := srcURL(apiBase, req.Workspace, req.RepoSlug, commit, req.ScopePath, true) + "?" + params.Encode()

	for urlStr != "" {
		resp, err := get(urlStr)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return err
		}
		list := new(SrcList)
		err = decodeResponse(resp, list)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		for _, entry := range list.Values {
			if !emit(entry) {
				return nil
			}
		}
		urlStr = list.Next
	}
	return nil
}

func readSrcFile(get func(string) (*http.Response, error), urlStr, path, commit string, emit func(SrcEntry) bool) error {
	resp, err := get(urlStr)
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Rows echo the requested path, because Postgres re-applies the path qualifier to every row.
	entry := SrcEntry{Path: path}
	if fullCommitHash.MatchString(commit) {
		entry.Commit.Hash = commit
	}

	if isJSONListing(resp.Header.Get("Content-Type"), body) {
		list := new(SrcList)
		if err := json.Unmarshal(body, list); err == nil && len(list.Values) > 0 {
			entry.Commit.Hash = list.Values[0].Commit.Hash
		}
		entry.Type = "commit_directory"
		emit(entry)
		return nil
	}

	size := int64(len(body))
	entry.Type = "commit_file"
	entry.Size = &size
	if utf8.Valid(body) && bytes.IndexByte(body, 0) < 0 {
		content := string(body)
		entry.Content = &content
	}
	emit(entry)
	return nil
}

// isJSONListing reports whether a src response is a directory listing. Files are served as
// text/plain or octet-stream whatever their extension, so a package.json is never mistaken for one.
func isJSONListing(contentType string, body []byte) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "application/json" {
		return false
	}
	var probe struct {
		Values *json.RawMessage `json:"values"`
	}
	return json.Unmarshal(body, &probe) == nil && probe.Values != nil
}

func srcURL(apiBase, workspace, repoSlug, commit, path string, directory bool) string {
	segments := []string{}
	for _, s := range strings.Split(strings.Trim(path, "/"), "/") {
		if s != "" {
			segments = append(segments, url.PathEscape(s))
		}
	}
	u := fmt.Sprintf("%s/repositories/%s/%s/src/%s/%s",
		strings.TrimRight(apiBase, "/"), url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(commit), strings.Join(segments, "/"))
	if directory && len(segments) > 0 {
		u += "/"
	}
	return u
}
