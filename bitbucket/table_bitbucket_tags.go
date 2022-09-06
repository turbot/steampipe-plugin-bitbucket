package bitbucket

import (
	"context"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
				Transform:   transform.FromQual("repository_full_name"),
			},
			{
				Name:        "self_link",
				Description: "A link to a resource related to this object.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource. It will be always \"tag\".",
				Type:        proto.ColumnType_STRING,
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

	for {
		response, err := client.Repositories.Repository.ListTags(opts)
		if err != nil {
			return nil, err
		}

		// Streaming result of current page
		for _, tag := range response.Tags {
			d.StreamListItem(ctx, tag)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Check for leftover pages (if any)
		if response.Next == "" {
			break
		}
		opts.PageNum = response.Page + 1
	}

	return nil, nil
}
