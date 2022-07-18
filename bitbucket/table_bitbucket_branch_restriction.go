package bitbucket

import (
	"context"
	"errors"
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableBitbucketBranchRestriction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_branch_restriction",
		Description: "Branch restrictions allow you to control the actions users can perform on a single branch, branch type, or branch pattern within a repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketBranchRestrictionsList,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "id",
				Description: "The unique ID if the branch restriction.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "self_link",
				Description: "The URL to the branch restriction.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "repository_full_name",
				Description: "The repository's full name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_full_name"),
			},

			// other fields
			{
				Name:        "branch_match_kind",
				Description: "The Branch match kind for the branch restriction.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "branch_type",
				Description: "The type of branch selected while creating the branch restriction.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "The type of restriction achieved using the branch restriction.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pattern",
				Description: "The branch name pattern specified while creating the branch restriction.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the branch operation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value associated to the kind for the branch restriction.",
				Type:        proto.ColumnType_INT,
			},

			// json fields
			{
				Name:        "groups",
				Description: "Details of the groups associated with the branch restriction.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "users",
				Description: "Details of the users associated with the branch restriction.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
		},
	}
}

func tableBitbucketBranchRestrictionsList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketBranchRestrictionsList")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repository_full_name should be in the format \"{workspace_slug}/{repo_slug}\"")
	}

	client := connect(ctx, d)

	opts := &bitbucket.BranchRestrictionsOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	response, err := client.Repositories.BranchRestrictions.Gets(opts)
	if err != nil {
		if isForbiddenError(err) {
			return nil, errors.New("Admin access to the repository is required in order to list the branch restrictions.")
		}
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketBranchRestrictionsList", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}
	branchRestrictionList := new(BranchRestrictionList)

	err = decodeJson(response, branchRestrictionList)
	if err != nil {
		return nil, err
	}

	for _, branchRestriction := range branchRestrictionList.BranchRestrictions {
		d.StreamListItem(ctx, branchRestriction)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

type BranchRestrictionList struct {
	ListResponse
	BranchRestrictions []BranchRestriction `json:"values,omitempty"`
}

type BranchRestriction struct {
	ID              int         `json:"id,omitempty"`
	Pattern         string      `json:"pattern,omitempty"`
	Kind            string      `json:"kind,omitempty"`
	Users           interface{} `json:"users,omitempty"`
	Links           interface{} `json:"links,omitempty"`
	Value           *int        `json:"value,omitempty"`
	BranchMatchKind string      `json:"branch_match_kind,omitempty"`
	Groups          interface{} `json:"groups,omitempty"`
	BranchType      string      `json:"branch_type,omitempty"`
	Type            string      `json:"type,omitempty"`
}
