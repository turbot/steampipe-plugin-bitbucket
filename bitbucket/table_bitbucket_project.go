package bitbucket

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// https://developer.atlassian.com/bitbucket/api/2/reference/resource/workspaces/%7Bworkspace%7D/projects
func tableBitbucketProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_project",
		Description: "A Bitbucket project. Projects are used by teams to organize repositories.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("workspace_slug"),
			Hydrate:    tableBitbucketProjectList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"workspace_slug", "key"}),
			Hydrate:    tableBitbucketProjectGet,
		},
		Columns: bitbucketProjectColumns(),
	}
}

func tableBitbucketProjectList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketProjectList")
	workspace := d.EqualsQuals["workspace_slug"].GetStringValue()
	if workspace == "" {
		return nil, nil
	}

	client := connect(ctx, d)
	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/projects", workspace)

	for {
		resp, err := client.HttpClient.Get(urlStr)
		if err != nil {
			if isNotFoundError(err) {
				return nil, nil
			}
			return nil, err
		}
		projectList := new(ProjectList)
		err = decodeResponse(resp, projectList)
		if err != nil {
			return nil, err
		}

		for _, project := range projectList.Projects {
			d.StreamListItem(ctx, project)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if projectList.Next == "" {
			return nil, nil
		}
		// update urlstring with the link of next page
		urlStr = projectList.Next
	}
}

func tableBitbucketProjectGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketProjectGet")
	workspace := d.EqualsQuals["workspace_slug"].GetStringValue()
	key := d.EqualsQuals["key"].GetStringValue()
	if workspace == "" || key == "" {
		return nil, nil
	}

	client := connect(ctx, d)
	urlStr := client.GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/projects/%s", workspace, key)

	resp, err := client.HttpClient.Get(urlStr)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	project := new(Project)
	err = decodeResponse(resp, project)
	if err != nil {
		return nil, err
	}

	return project, nil

}
