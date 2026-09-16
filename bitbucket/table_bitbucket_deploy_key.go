package bitbucket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-get
func tableBitbucketDeployKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_deploy_key",
		Description: "Access keys configured on a Bitbucket repository. A deploy key grants access to the repository independently of any user account.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketDeployKeyList,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The key's id, unique within the repository.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "repository_full_name",
				Description: "Full name of the repository the key grants access to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_full_name"),
			},
			{
				Name:        "label",
				Description: "Label given to the key when it was added.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "comment",
				Description: "Comment parsed from the public key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_on",
				Description: "Time when the key was added to the repository. This is what dates a key for rotation; Bitbucket never reports a key's age any other way.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_used",
				Description: "Time the key was last used, or null if it has never been used.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "key",
				Description: "The public key.",
				Type:        proto.ColumnType_STRING,
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
				Transform:   transform.FromField("Label"),
			},
		},
	}
}

func tableBitbucketDeployKeyList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketDeployKeyList")

	repoFullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repoSlug := parseRepoFullName(repoFullName)
	if owner == "" || repoSlug == "" {
		return nil, nil
	}

	client := connect(ctx, d)
	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/repositories/%s/%s/deploy-keys", owner, repoSlug)

	err := page(ctx, d, urlStr, func(resp *http.Response) (string, bool, error) {
		keyList := new(DeployKeyList)
		if err := decodeResponse(resp, keyList); err != nil {
			return "", false, err
		}
		for _, item := range keyList.Keys {
			d.StreamListItem(ctx, item)
			if d.RowsRemaining(ctx) == 0 {
				return "", false, nil
			}
		}
		return keyList.Next, true, nil
	})

	return nil, err
}

type DeployKeyList struct {
	ListResponse
	Keys []DeployKey `json:"values,omitempty"`
}

type DeployKey struct {
	Comment   string     `json:"comment"`
	CreatedOn *time.Time `json:"created_on"`
	Id        int64      `json:"id"`
	Key       string     `json:"key"`
	Label     string     `json:"label"`
	LastUsed  *time.Time `json:"last_used"`
	Type      string     `json:"type"`
}
