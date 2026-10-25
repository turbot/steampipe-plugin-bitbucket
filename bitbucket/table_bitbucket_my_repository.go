package bitbucket

import (
	"context"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func tableBitbucketDefaultReviewersList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("tableBitbucketDefaultReviewersList")

	// `bitbucket_repository` streams Repo and `bitbucket_my_repository` streams
	// bitbucket.Repository, and both share this column set. Owner is derived from full_name rather
	// than the owner map, which carries no display_name for a workspace-owned repository.
	var fullName, uuid string
	switch data := h.Item.(type) {
	case Repo:
		fullName, uuid = data.Full_name, data.Uuid
	case bitbucket.Repository:
		fullName, uuid = data.Full_name, data.Uuid
	default:
		return nil, nil
	}
	owner, repoSlug := parseRepoFullName(fullName)
	if owner == "" || repoSlug == "" {
		return nil, nil
	}

	client := connect(ctx, d)

	opts := &bitbucket.RepositoryOptions{
		Owner:    owner,
		Uuid:     uuid,
		RepoSlug: repoSlug,
	}

	response, err := client.Repositories.Repository.ListDefaultReviewers(opts)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketDefaultReviewersList", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}
	return response.DefaultReviewers, nil
}
