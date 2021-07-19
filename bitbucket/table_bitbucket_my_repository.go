package bitbucket

import (
	"context"
	"fmt"

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
	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/repositories/%s", owner)

	for {
		resp, err := client.HttpClient.Get(urlStr)
		if err != nil {
			if isForbiddenError(err) {
				return nil, nil
			}
			if isNotFoundError(err) {
				return nil, nil
			}
			plugin.Logger(ctx).Error("tableBitbucketRepositoryList", "Error", err)
			return nil, err
		}

		repositoryList := new(RepositoryList)
		err = decodeResponse(resp, repositoryList)
		if err != nil {
			return nil, err
		}
		for _, repository := range repositoryList.Repositories {
			d.StreamListItem(ctx, repository)
		}
		if repositoryList.Next == "" {
			return nil, nil
		}
		// update urlstring with the link of next page
		urlStr = repositoryList.Next
	}
}
