package bitbucket

import (
	"context"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func bitBucketRepositoryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "name",
			Description: "Name",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "uuid",
			Description: "Uuid",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Uuid"),
		},
		{
			Name:        "slug",
			Description: "Slug",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "full_name",
			Description: "Full_name",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Full_name"),
		},
		{
			Name:        "description",
			Description: "Description",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "fork_policy",
			Description: "ForkPolicy",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "language",
			Description: "Language",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "is_private",
			Description: "Is_private",
			Type:        proto.ColumnType_BOOL,
			Transform:   transform.FromField("is_private"),
		},
		{
			Name:        "has_issues",
			Description: "Has_issues",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("has_issues"),
		},
		{
			Name:        "mainbranch",
			Description: "Mainbranch",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "type",
			Description: "Type",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "owner",
			Description: "Owner",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "links",
			Description: "Links",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "parent",
			Description: "Parent",
			Type:        proto.ColumnType_JSON,
		},
	}
}

func tableBitBucketRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_repository",
		Description: "BitBucket Repositories that you are associated with.  BitBucket Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("full_name"),
			Hydrate:    tableBitbucketRepositoryList,
		},
		Columns: bitBucketRepositoryColumns(),
	}
}

func tableBitbucketRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	repoFullName := d.KeyColumnQuals["full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	client := connect(ctx, d)
	opts := &bitbucket.RepositoryOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	repo, err := client.Repositories.Repository.Get(opts)
	if err != nil {
		return nil, err
	}

	if repo != nil {
		d.StreamListItem(ctx, repo)
	}

	return nil, nil
}
