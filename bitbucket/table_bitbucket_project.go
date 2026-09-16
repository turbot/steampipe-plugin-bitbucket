package bitbucket

import (
	"context"
	"fmt"
	"net/http"

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

	urlStr := connect(ctx, d).GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/projects", workspace)

	err := page(ctx, d, urlStr, func(resp *http.Response) (string, bool, error) {
		projectList := new(ProjectList)
		if err := decodeResponse(resp, projectList); err != nil {
			return "", false, err
		}
		for _, project := range projectList.Projects {
			d.StreamListItem(ctx, project)
			if d.RowsRemaining(ctx) == 0 {
				return "", false, nil
			}
		}
		return projectList.Next, true, nil
	})

	return nil, err
}

func tableBitbucketProjectGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketProjectGet")
	workspace := d.EqualsQuals["workspace_slug"].GetStringValue()
	key := d.EqualsQuals["key"].GetStringValue()
	if workspace == "" || key == "" {
		return nil, nil
	}

	urlStr := connect(ctx, d).GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/projects/%s", workspace, key)

	resp, err := authGet(ctx, d, urlStr)
	if err != nil {
		if isNotFoundError(err) || isForbiddenError(err) {
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
