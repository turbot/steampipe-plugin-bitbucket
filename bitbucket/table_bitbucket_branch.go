package bitbucket

import (
	"context"
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableBitbucketBranch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_branch",
		Description: "BitBucket tags associated with the repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketBranchList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "name"}),
			Hydrate:    tableBitbucketBranchGet,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "name",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_full_name",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_merge_strategy",
				Description: "",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "heads",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "links",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "merge_strategies",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

type branchList struct {
	bitbucket.RepositoryBranch
	RepositoryFullName string
}

func tableBitbucketBranchList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	client := connect(ctx, d)
	opts := &bitbucket.RepositoryBranchOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	branches, err := client.Repositories.Repository.ListBranches(opts)
	if err != nil {
		return nil, err
	}
	if branches != nil {
		for _, branch := range branches.Branches {
			d.StreamListItem(ctx, branchList{branch, repoFullName})
		}
	}

	return nil, nil
}

func tableBitbucketBranchGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketBranchGet")

	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	name := d.KeyColumnQuals["name"].GetStringValue()

	if repoFullName == "" || name == "" {
		return nil, nil
	}
	owner, repoName := parseRepoFullName(repoFullName)

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repository_full_name should be in the format \"{workspace_slug}/{repo_slug}\"")
	}

	client := connect(ctx, d)

	opts := &bitbucket.RepositoryBranchOptions{
		Owner:      owner,
		RepoSlug:   repoName,
		BranchName: name,
	}

	response, err := client.Repositories.Repository.GetBranch(opts)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketBranchGet", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	return branchList{*response, repoFullName}, nil
}
