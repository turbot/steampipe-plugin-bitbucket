package bitbucket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableBitbucketRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_repository",
		Description: "BitBucket repositorie that you are associated with. BitBucket repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"full_name", "workspace_slug"}),
			Hydrate:    tableBitbucketRepositoryList,
		},
		Columns: bitBucketRepositoryColumns(),
	}
}

// Qualifying on `workspace_slug` lists a workspace's repositories; qualifying on `full_name`
// fetches one. `bitbucket_my_repository` cannot substitute for the former: it hangs off
// `bitbucket_my_workspace`, whose `/2.0/workspaces` endpoint Atlassian has withdrawn (404), so
// without this there is no way to list the repositories of a workspace you can name.
func tableBitbucketRepositoryList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	if workspace := d.EqualsQuals["workspace_slug"].GetStringValue(); workspace != "" {
		return listWorkspaceRepositories(ctx, d, workspace)
	}

	repoFullName := d.EqualsQuals["full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)
	if owner == "" || repoName == "" {
		return nil, nil
	}

	urlStr := connect(ctx, d).GetApiBaseURL() + fmt.Sprintf("/repositories/%s/%s", owner, repoName)

	resp, err := authGet(ctx, d, urlStr)
	if err != nil {
		if isNotFoundError(err) || isForbiddenError(err) {
			return nil, nil
		}
		return nil, err
	}

	repo := new(Repo)
	if err := decodeResponse(resp, repo); err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, *repo)

	return nil, nil
}

func listWorkspaceRepositories(ctx context.Context, d *plugin.QueryData, workspace string) (interface{}, error) {
	urlStr := connect(ctx, d).GetApiBaseURL() + fmt.Sprintf("/repositories/%s", workspace)

	err := page(ctx, d, urlStr, func(resp *http.Response) (string, bool, error) {
		repoList := new(RepoList)
		if err := decodeResponse(resp, repoList); err != nil {
			return "", false, err
		}
		for _, repo := range repoList.Repositories {
			d.StreamListItem(ctx, repo)
			if d.RowsRemaining(ctx) == 0 {
				return "", false, nil
			}
		}
		return repoList.Next, true, nil
	})

	return nil, err
}

// Field names match `bitbucket.Repository` so the existing column transforms are unchanged. The
// API carries more than that struct models — `created_on`, `updated_on`, `size` and `scm` among
// them — and those are what tells a stale repository from a live one.
type Repo struct {
	CreatedOn             *time.Time             `json:"created_on"`
	Description           string                 `json:"description"`
	EnforcedSignedCommits *bool                  `json:"enforced_signed_commits"`
	ForkPolicy            string                 `json:"fork_policy"`
	Full_name             string                 `json:"full_name"`
	Has_issues            bool                   `json:"has_issues"`
	Is_private            bool                   `json:"is_private"`
	Language              string                 `json:"language"`
	Links                 map[string]interface{} `json:"links"`
	Mainbranch            map[string]interface{} `json:"mainbranch"`
	Name                  string                 `json:"name"`
	Owner                 map[string]interface{} `json:"owner"`
	Parent                map[string]interface{} `json:"parent"`
	Project               RepoProject            `json:"project"`
	Scm                   string                 `json:"scm"`
	Size                  int64                  `json:"size"`
	Slug                  string                 `json:"slug"`
	Type                  string                 `json:"type"`
	UpdatedOn             *time.Time             `json:"updated_on"`
	Uuid                  string                 `json:"uuid"`
	Website               string                 `json:"website"`
	Workspace             RepoWorkspace          `json:"workspace"`
}

type RepoProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

type RepoWorkspace struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Uuid string `json:"uuid"`
}

type RepoList struct {
	ListResponse
	Repositories []Repo `json:"values,omitempty"`
}

func bitBucketRepositoryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "name",
			Description: "The name of repository.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "uuid",
			Description: "The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Uuid"),
		},
		{
			Name:        "slug",
			Description: "A repository slug is a URL-friendly version of a repository name, automatically generated by Bitbucket for use in the URL.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "full_name",
			Description: "The concatenation of the repository owner's username and the slugified name, e.g. \"turbot/steampipe-plugin-bitbucket\". This is the same string used in Bitbucket URLs.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Full_name"),
		},
		{
			Name:        "workspace_slug",
			Description: "Slug of the workspace this repository belongs to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Workspace.Slug"),
		},
		{
			Name:        "description",
			Description: "Description of the repository.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "created_on",
			Description: "Time when the repository was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_on",
			Description: "Time of the most recent change to the repository, including pushes.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "size",
			Description: "Size of the repository in bytes.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "scm",
			Description: "The source control system the repository uses, e.g. \"git\".",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "website",
			Description: "Website configured for the repository.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "enforced_signed_commits",
			Description: "True if the repository rejects unsigned commits. Null when Bitbucket does not report it.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "fork_policy",
			Description: "Controls the rules for forking this repository. \"allow_forks\": unrestricted forking, \"no_public_forks\": restrict forking to private forks (forks cannot be made public later) and \"no_forks\": deny all forking",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "language",
			Description: "The type of markup language the raw content is to be interpreted in.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "is_private",
			Description: "Indicates whether the repository is publicly accessible, or whether it is private to the team and consequently only visible to team members.",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("Is_private"),
		},
		{
			Name:        "has_issues",
			Description: "To initialize or disable the new repo's issue tracker",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Has_issues"),
		},
		{
			Name:        "owner_account_id",
			Description: "Jira account id of the owner.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Owner.account_id"),
		},
		{
			Name:        "owner_display_name",
			Description: "Display name of the owner the repository.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Owner.display_name"),
		},
		{
			Name:        "owner_type",
			Description: "Type of the owner of the repository. Can be a user or team.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Owner.type"),
		},
		{
			Name:        "owner_uuid",
			Description: "Bitbucket UUID of the owner.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Owner.uuid"),
		},
		{
			Name:        "project_name",
			Description: "Name of the project this repository belongs to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Project.Name"),
		},
		{
			Name:        "project_key",
			Description: "Key of the project this repository belongs to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Project.Key"),
		},
		{
			Name:        "project_uuid",
			Description: "UUID of the project this repository belongs to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Project.Uuid"),
		},
		{
			Name:        "self_link",
			Description: "Self link to this repository.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Links.self.href"),
		},
		{
			Name:        "mainbranch",
			Description: "Details of the main branch of the repository.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "default_reviewers",
			Description: "Details of the default reviewers of the repository.",
			Type:        proto.ColumnType_JSON,
			Hydrate:     tableBitbucketDefaultReviewersList,
			Transform:   transform.FromValue(),
		},
		// Standard columns
		{
			Name:        "title",
			Description: ColumnDescriptionTitle,
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Name"),
		},
	}
}
