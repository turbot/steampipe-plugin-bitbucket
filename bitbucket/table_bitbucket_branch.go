package bitbucket

import (
	"context"
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableBitbucketBranch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_branch",
		Description: "BitBucket Branch.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketBranchesList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "name"}),
			Hydrate:    tableBitbucketBranchGet,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "name",
				Description: "The name of the branch.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_full_name",
				Description: "The repository full name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_full_name"),
			},
			{
				Name:        "type",
				Description: "The branch type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_merge_strategy",
				Description: "Branch default merge strategy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Default_Merge_Strategy"),
			},

			// json fields
			{
				Name:        "heads",
				Description: "Branch head details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "links",
				Description: "Branch link details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "merge_strategies",
				Description: "Branch merge strategies.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Merge_Strategies"),
			},
			{
				Name:        "target",
				Description: "Branch target details.",
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

func tableBitbucketBranchesList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	client := connect(ctx, d)
	opts := &bitbucket.RepositoryBranchOptions{
		Owner:    owner,
		RepoSlug: repoName,
		Pagelen:  100,
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if int(*limit) < opts.Pagelen {
			if *limit < 1 {
				opts.Pagelen = 1
			} else {
				opts.Pagelen = int(*limit)
			}
		}
	}

	pagesLeft := true

	for pagesLeft {
		branches, err := client.Repositories.Repository.ListBranches(opts)
		if err != nil {
			return nil, err
		}

		if branches != nil {
			for _, branch := range branches.Branches {
				d.StreamListItem(ctx, branch)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}

		if branches.Next == "" {
			pagesLeft = false
		} else {
			opts.PageNum = branches.Page + 1
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

	return response, nil
}
