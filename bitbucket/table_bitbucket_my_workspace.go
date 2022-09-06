package bitbucket

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableBitbucketMyWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_my_workspace",
		Description: "Workspace is where you will create repositories, collaborate on your code, and organize different streams of work in your Bitbucket Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: tableBitbucketMyWorkspaceList,
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

func tableBitbucketMyWorkspaceList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketWorkspaceList")
	client := connect(ctx, d)

	resp, err := client.Workspaces.List()
	if err != nil {
		return nil, err
	}

	for _, workspace := range resp.Workspaces {
		d.StreamListItem(ctx, workspace)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
