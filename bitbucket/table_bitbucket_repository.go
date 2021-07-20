package bitbucket

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableBitbucketRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_repository",
		Description: "BitBucket repositorie that you are associated with. BitBucket repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("full_name"),
			Hydrate:    tableBitbucketRepositoryList,
		},
		Columns: bitBucketRepositoryColumns(),
	}
}

func tableBitbucketRepositoryList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	repoFullName := d.KeyColumnQuals["full_name"].GetStringValue()
	if repoFullName == "" {
		return nil, nil
	}
	owner, repoName := parseRepoFullName(repoFullName)
	client := connect(ctx, d)

	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/repositories/%s/%s", owner, repoName)
	resp, err := client.HttpClient.Get(urlStr)
	if err != nil {
		if isForbiddenError(err) {
			return nil, nil
		}
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketRepositoryList", "Error", err)
		return nil, err
	}
	repository := new(Repository)
	err = decodeResponse(resp, repository)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			return nil, nil
		}
		return nil, err
	}
	d.StreamListItem(ctx, repository)

	return nil, nil
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
			Transform:   transform.FromField("UUID"),
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
			Transform:   transform.FromField("FullName"),
		},
		{
			Name:        "created",
			Description: "Timestamp when the repository was created.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("CreatedOn"),
		},
		{
			Name:        "description",
			Description: "Description of the repository.",
			Type:        proto.ColumnType_STRING,
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
		},
		{
			Name:        "has_issues",
			Description: "To initialize or disable the new repo's issue tracker",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "has_wiki",
			Description: "Indicates whether the repository is having a Wiki or not.",
			Type:        proto.ColumnType_BOOL,
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
			Transform:   transform.FromField("Project.name"),
		},
		{
			Name:        "project_key",
			Description: "Key of the project this repository belongs to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Project.key"),
		},
		{
			Name:        "project_uuid",
			Description: "UUID of the project this repository belongs to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Project.uuid"),
		},
		{
			Name:        "self_link",
			Description: "Self link to this repository.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Links.self.href"),
		},
		{
			Name:        "website",
			Description: "Self link to this repository.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "updated",
			Description: "Timestamp when the repository was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
			Transform:   transform.FromField("UpdatedOn"),
		},
		{
			Name:        "mainbranch",
			Description: "Details of the main branch of the repository.",
			Type:        proto.ColumnType_JSON,
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

//// Custom Structs

type RepositoryList struct {
	ListResponse
	Repositories []Repository `json:"values,omitempty"`
}

type Repository struct {
	Website     string                 `json:"website,omitempty"`
	HasWiki     bool                   `json:"has_wiki,omitempty"`
	UUID        string                 `json:"uuid,omitempty"`
	Links       map[string]interface{} `json:"links,omitempty"`
	ForkPolicy  string                 `json:"fork_policy,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Project     map[string]interface{} `json:"project,omitempty"`
	Language    string                 `json:"language,omitempty"`
	CreatedOn   *time.Time             `json:"created_on,omitempty"`
	Mainbranch  map[string]interface{} `json:"mainbranch,omitempty"`
	Workspace   map[string]interface{} `json:"workspace,omitempty"`
	HasIssues   bool                   `json:"has_issues,omitempty"`
	Owner       map[string]interface{} `json:"owner,omitempty"`
	UpdatedOn   *time.Time             `json:"updated_on,omitempty"`
	Slug        string                 `json:"slug,omitempty"`
	IsPrivate   bool                   `json:"is_private,omitempty"`
	Description string                 `json:"description,omitempty"`
	FullName    string                 `json:"full_name,omitempty"`
}
