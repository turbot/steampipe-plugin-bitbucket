package bitbucket

import (
	"context"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableBitbucketMyRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_my_repository",
		Description: "BitBucket repositories that you are associated with. BitBucket Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			ParentHydrate: tableBitbucketMyWorkspaceList,
			Hydrate:       tableBitbucketMyRepositoryList,
		},
		Columns: bitBucketRepositoryColumns(),
	}
}

func tableBitbucketMyRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	owner := h.Item.(bitbucket.Workspace).Slug
	client := connect(ctx, d)

	repos, err := client.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
		Owner: owner,
	})

	if err != nil {
		return nil, err
	}

	for _, repo := range repos.Items {
		d.StreamListItem(ctx, repo)
	}

	return nil, nil
}
