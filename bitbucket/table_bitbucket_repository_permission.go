package bitbucket

import (
	"context"
	"fmt"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-permissions-repositories-get
func tableBitbucketRepositoryPermission(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_repository_permission",
		Description: "Explicit user permissions on the repositories of a Bitbucket workspace.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"workspace_slug", "repository_full_name"}),
			Hydrate:    tableBitbucketRepositoryPermissionList,
		},
		Columns: []*plugin.Column{
			{
				Name:        "workspace_slug",
				Description: "Slug of the workspace the repository belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_full_name",
				Description: "Full name of the repository, e.g. \"my-workspace/my-repo\".",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.FullName"),
			},
			{
				Name:        "repository_uuid",
				Description: "The repository's immutable id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.Uuid"),
			},
			{
				Name:        "permission",
				Description: "The user's permission on the repository: admin, write or read.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_uuid",
				Description: "The user's immutable id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.UUID"),
			},
			{
				Name:        "user_account_id",
				Description: "Atlassian account id of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.AccountId"),
			},
			{
				Name:        "user_display_name",
				Description: "Display name of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.DisplayName"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.DisplayName"),
			},
		},
	}
}

// Bitbucket exposes repository permissions only per workspace, so a `repository_full_name` qual is
// served by filtering that listing rather than by a narrower endpoint.
func tableBitbucketRepositoryPermissionList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketRepositoryPermissionList")

	repoFullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	workspaceSlug := d.EqualsQuals["workspace_slug"].GetStringValue()
	if workspaceSlug == "" {
		workspaceSlug, _ = parseRepoFullName(repoFullName)
	}
	if workspaceSlug == "" {
		return nil, nil
	}

	client := connect(ctx, d)
	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/permissions/repositories", workspaceSlug)
	if repoFullName != "" {
		_, repoSlug := parseRepoFullName(repoFullName)
		if repoSlug == "" {
			return nil, nil
		}
		urlStr = client.GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/permissions/repositories/%s", workspaceSlug, repoSlug)
	}

	err := page(ctx, d, urlStr, func(resp *http.Response) (string, bool, error) {
		permissionList := new(RepositoryPermissionList)
		if err := decodeResponse(resp, permissionList); err != nil {
			return "", false, err
		}
		for _, item := range permissionList.Permissions {
			// The response echoes the repository but not the workspace it was listed under.
			item.WorkspaceSlug = workspaceSlug
			d.StreamListItem(ctx, item)
			if d.RowsRemaining(ctx) == 0 {
				return "", false, nil
			}
		}
		return permissionList.Next, true, nil
	})

	return nil, err
}

type RepositoryPermissionList struct {
	ListResponse
	Permissions []RepositoryPermission `json:"values,omitempty"`
}

type RepositoryPermission struct {
	Type          string            `json:"type"`
	Permission    string            `json:"permission"`
	User          User              `json:"user"`
	Repository    PermissionRepoRef `json:"repository"`
	WorkspaceSlug string            `json:"-"`
}

type PermissionRepoRef struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Uuid     string `json:"uuid"`
}
