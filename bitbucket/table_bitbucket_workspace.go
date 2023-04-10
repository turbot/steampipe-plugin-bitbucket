package bitbucket

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableBitbucketWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_workspace",
		Description: "Workspace is where you will create repositories, collaborate on your code, and organize different streams of work in your Bitbucket Cloud account.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("slug"),
			Hydrate:    tableBitbucketWorkspaceGet,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "slug",
				Description: "The short label that identifies this workspace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uuid",
				Description: "The workspace's immutable id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "is_private",
				Description: "Indicates whether the workspace is publicly accessible, or whether it is private to the members and consequently only visible to members. Note that private workspaces cannot contain public repositories.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Is_Private"),
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource.",
				Type:        proto.ColumnType_STRING,
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

func tableBitbucketWorkspaceGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketWorkspaceGet")

	slug := d.EqualsQuals["slug"].GetStringValue()
	if slug == "" {
		return nil, nil
	}

	client := connect(ctx, d)

	workspace, err := client.Workspaces.Get(slug)
	if err != nil {
		if isNotFoundError(err) || isForbiddenError(err) {
			return nil, nil
		}
		return nil, err
	}

	if workspace == nil {
		return nil, nil
	}

	d.StreamListItem(ctx, workspace)

	return nil, nil
}
