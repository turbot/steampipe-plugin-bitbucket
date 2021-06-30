package bitbucket

import (
	"context"
	"time"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// https://developer.atlassian.com/bitbucket/api/2/reference/resource/workspaces/%7Bworkspace%7D/members
func tableBitbucketWorkspaceMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_workspace_member",
		Description: "A Bitbucket project. Projects are used by teams to organize repositories.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("workspace_slug"),
			Hydrate:    tableBitbucketWorkspaceMemberList,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "display_name",
				Description: "Display name of the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.DisplayName"),
			},
			{
				Name:        "uuid",
				Description: "The member's immutable id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UUID"),
			},
			{
				Name:        "account_id",
				Description: "Account id of the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.AccountId"),
			},
			{
				Name:        "self_link",
				Description: "Self link to the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "member_type",
				Description: "Type of the member.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.Type"),
			},
			{
				Name:        "workspace_slug",
				Description: "Slug of the workspace to which this member belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Workspace.Slug"),
			},
		},
	}
}

func tableBitbucketWorkspaceMemberList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketWorkspaceMemberList")

	workspaceSlug := d.KeyColumnQuals["workspace_slug"].GetStringValue()
	client := connect(ctx, d)

	response, err := client.Workspaces.Members(workspaceSlug)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	memberList := new(MemberList)
	err = decodeJson(response, memberList)
	if err != nil {
		return nil, err
	}

	for _, item := range memberList.Members {
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}

type MemberList struct {
	ListResponse
	Members []Member `json:"values,omitempty"`
}

type Member struct {
	Type      string                 `json:"type"`
	Links     map[string]interface{} `json:"links"`
	User      User                   `json:"user"`
	Workspace bitbucket.Workspace    `json:"workspace"`
}

type User struct {
	AccountId     string                 `json:"account_id,omitempty"`
	AccountStatus string                 `json:"account_status,omitempty"`
	CreatedOn     *time.Time             `json:"created_on,omitempty"`
	DisplayName   string                 `json:"display_name,omitempty"`
	Has2faEnabled bool                   `json:"has_2fa_enabled,omitempty"`
	IsStaff       bool                   `json:"is_staff,omitempty"`
	Nickname      string                 `json:"nickname,omitempty"`
	Type          string                 `json:"type,omitempty"`
	UUID          string                 `json:"uuid,omitempty"`
	Website       string                 `json:"website,omitempty"`
	Links         map[string]interface{} `json:"links,omitempty"`
}
