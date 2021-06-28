package bitbucket

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableBitBucketWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_workspace",
		Description: "Workspace is where you will create repositories, collaborate on your code, and organize different streams of work in your Bitbucket Cloud account.",
		List: &plugin.ListConfig{
			Hydrate: tableBitbucketWorkspaceList,
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
		},
	}
}

func tableBitbucketWorkspaceList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketWorkspaceList")
	client := connect(ctx, d)

	resp, err := client.Workspaces.List()
	if err != nil {
		return nil, err
	}

	for _, workspace := range resp.Workspaces {
		d.StreamListItem(ctx, workspace)
	}

	return nil, nil
}
