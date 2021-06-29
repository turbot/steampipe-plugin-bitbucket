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
			"bitbucket_issue":            tableBitbucketIssue(ctx),
			"bitbucket_my_repository":    tableBitbucketMyRepository(ctx),
			"bitbucket_my_workspace":     tableBitbucketMyWorkspace(ctx),
			"bitbucket_project":          tableBitbucketProject(ctx),
			"bitbucket_repository":       tableBitbucketRepository(ctx),
			"bitbucket_snippet":          tableBitbucketSnippet(ctx),
			"bitbucket_workspace_member": tableBitbucketWorkspaceMember(ctx),
		},
	}
	return p
}
