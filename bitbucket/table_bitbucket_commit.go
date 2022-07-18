package bitbucket

import (
	"context"
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableBitbucketCommit(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_commit",
		Description: "A Bitbucket commit is an operation which sends the latest changes of the source code to the repository, making these changes part of the head revision of the repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketCommitsList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "hash"}),
			Hydrate:    tableBitbucketCommitGet,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "hash",
				Description: "The hash id of the commit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "message",
				Description: "The message of the commit.",
				Type:        proto.ColumnType_STRING,
			},

			// other fields
			{
				Name:        "author_account_id",
				Description: "The account id of the author of the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.user.account_id"),
			},
			{
				Name:        "author_display_name",
				Description: "The display name of the author of the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.user.display_name"),
			},
			{
				Name:        "author_self_link",
				Description: "The self link of the author of the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.user.links.self.href"),
			},
			{
				Name:        "author_type",
				Description: "The type of the author of the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.user.type"),
			},
			{
				Name:        "author_uuid",
				Description: "The UUID of the author of the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.user.uuid"),
			},
			{
				Name:        "repository_full_name",
				Description: "The repository's full name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.full_name"),
			},
			{
				Name:        "repository_name",
				Description: "The repository's name (slug).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.name"),
			},
			{
				Name:        "repository_self_link",
				Description: "The self link of the repository that contains the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.links.self.href"),
			},
			{
				Name:        "repository_uuid",
				Description: "The UUID of the repository that contains the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.uuid"),
			},
			{
				Name:        "self_link",
				Description: "The self link of the commit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "parents",
				Description: "Details of parent commit (if any).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "summary",
				Description: "A brief summary of the commit.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

func tableBitbucketCommitsList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketCommitsList")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repository_full_name should be in the format \"{workspace_slug}/{repo_slug}\"")
	}

	client := connect(ctx, d)

	opts := &bitbucket.CommitsOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	response, err := client.Repositories.Commits.GetCommits(opts)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketCommitsList", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}
	commitList := new(CommitList)
	err = decodeJson(response, commitList)
	if err != nil {
		return nil, err
	}

	for _, commit := range commitList.Commits {
		d.StreamListItem(ctx, commit)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func tableBitbucketCommitGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketCommitGet")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	revisionID := d.KeyColumnQuals["hash"].GetStringValue()

	if repoFullName == "" {
		return nil, nil
	}
	owner, repoName := parseRepoFullName(repoFullName)

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repository_full_name should be in the format \"{workspace_slug}/{repo_slug}\"")
	}

	if revisionID == "" {
		return nil, nil
	}
	client := connect(ctx, d)

	opts := &bitbucket.CommitsOptions{
		Owner:    owner,
		RepoSlug: repoName,
		Revision: revisionID,
	}

	response, err := client.Repositories.Commits.GetCommit(opts)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketCommitGet", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	commit := new(Commit)
	err = decodeJson(response, commit)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

type CommitList struct {
	ListResponse
	Commits []Commit `json:"values,omitempty"`
}

type Commit struct {
	Hash       string                 `json:"hash,omitempty"`
	Repository map[string]interface{} `json:"repository,omitempty"`
	Links      map[string]interface{} `json:"links,omitempty"`
	Author     map[string]interface{} `json:"author,omitempty"`
	Summary    map[string]interface{} `json:"summary,omitempty"`
	Parents    interface{}            `json:"parents,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Type       string                 `json:"type,omitempty"`
}
