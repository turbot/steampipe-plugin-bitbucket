package bitbucket

import (
	"context"
	"fmt"
	"time"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// https://developer.atlassian.com/bitbucket/api/2/reference/resource/workspaces/%7Bworkspace%7D/projects
func tableBitBucketProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_project",
		Description: "A Bitbucket project. Projects are used by teams to organize repositories.",
		List: &plugin.ListConfig{
			ParentHydrate: tableBitbucketWorkspaceList,
			Hydrate:       tableBitbucketProjectList,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "name",
				Description: "The name of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uuid",
				Description: "The project's immutable id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "key",
				Description: "The project's key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "workspace",
				Description: "Slug name of the workspace to which this project belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Workspace.Slug"),
			},

			// other fields
			{
				Name:        "created",
				Description: "Timestamp when project was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_private",
				Description: "Indicates whether the project is publicly accessible, or whether it is private to the team and consequently only visible to team members. Note that private projects cannot contain public repositories.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "owner_display_name",
				Description: "Display name of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.display_name"),
			},
			{
				Name:        "owner_type",
				Description: "Type of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.type"),
			},
			{
				Name:        "owner_uuid",
				Description: "UUID of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.uuid"),
			},
			{
				Name:        "self_link",
				Description: "A self link to this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated",
				Description: "Timestamp when project was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
		},
	}
}

func tableBitbucketProjectList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketProjectList")
	workspace := h.Item.(bitbucket.Workspace)
	client := connect(ctx, d)

	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/projects", workspace.Slug)

	for {
		resp, err := client.HttpClient.Get(urlStr)
		if err != nil {
			return nil, err
		}
		projectList := new(ProjectList)
		err = decodeResponse(resp, projectList)
		if err != nil {
			return nil, err
		}

		for _, project := range projectList.Projects {
			d.StreamListItem(ctx, project)
		}

		if projectList.Next == "" {
			return nil, nil
		}
		// update urlstring with the link of next page
		urlStr = projectList.Next
	}
}

type ProjectList struct {
	Page     int       `json:"page,omitempty"`
	Pagelen  int       `json:"pagelen,omitempty"`
	MaxDepth int       `json:"maxDepth,omitempty"`
	Size     int       `json:"size,omitempty"`
	Next     string    `json:"next,omitempty"`
	Previous string    `json:"previous,omitempty"`
	Projects []Project `json:"values,omitempty"`
}

type Project struct {
	Created     *time.Time             `json:"created_on"`
	Description string                 `json:"description"`
	IsPrivate   bool                   `json:"is_private"`
	Key         string                 `json:"key"`
	Links       map[string]interface{} `json:"links"`
	Name        string                 `json:"name"`
	Owner       map[string]interface{} `json:"owner"`
	Type        string                 `json:"type"`
	UUID        string                 `json:"uuid"`
	Updated     *time.Time             `json:"updated_on"`
	Workspace   bitbucket.Workspace    `json:"workspace"`
}
