package bitbucket

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// Plugin returns this plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-bitbucket",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromCamel(),
		TableMap: map[string]*plugin.Table{
			"bitbucket_my_repository":    tableBitBucketMyRepository(ctx),
			"bitbucket_project":          tableBitBucketProject(ctx),
			"bitbucket_repository":       tableBitBucketRepository(ctx),
			"bitbucket_workspace":        tableBitBucketWorkspace(ctx),
			"bitbucket_workspace_member": tableBitBucketWorkspaceMember(ctx),
		},
	}
	return p
}
