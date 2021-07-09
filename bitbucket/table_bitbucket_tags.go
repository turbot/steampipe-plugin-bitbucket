package bitbucket

import (
	"context"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableBitbucketTag(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_tag",
		Description: "Bitbucket repository tags.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketTagsList,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the ref.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_full_name",
				Description: "The concatenation of the repository owner's username and the slugified name, e.g. \"turbot/steampipe-plugin-bitbucket\". This is the same string used in Bitbucket URLs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource. It will be always \"tag\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "links",
				Description: "A link to a resource related to this object.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target",
				Description: "Specifies details of target of the tag.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "heads",
				Description: "The repository's full name.",
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

type repositoryTags = struct {
	bitbucket.RepositoryTag
	RepositoryFullName string
}

//// LIST FUNCTION

func tableBitbucketTagsList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketTagsList")

	client := connect(ctx, d)
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	// Build params
	opts := &bitbucket.RepositoryTagOptions{
		Owner:    owner,
		RepoSlug: repoName,
		Pagelen: 100,
	}

	for {
		response, err := client.Repositories.Repository.ListTags(opts)
		if err != nil {
			return nil, err
		}

		// Streaming result of current page
		for _, tag := range response.Tags {
			d.StreamListItem(ctx, repositoryTags{tag, repoFullName})
		}

		// Check for leftover pages (if any)
		if response.Next == "" {
			break
		}
		opts.PageNum = response.Page + 1
	}

	return nil, nil
}
