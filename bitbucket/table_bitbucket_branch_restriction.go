package bitbucket

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableBitbucketBranchRestriction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_branch_restriction",
		Description: "TODO",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketBranchRestrictionsList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    tableBitbucketBranchRestrictionGet,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "id",
				Description: "The unique ID if the branch restriction.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			// {
			// 	Name:        "self_link",
			// 	Description: "TODO.",
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("Links.self.href"),
			// },
			{
				Name:        "repository_full_name",
				Description: "The repository's full name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_full_name"),
			},

			// other fields
			{
				Name:        "branch_match_kind",
				Description: "TODO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "branch_type",
				Description: "TODO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "TODO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pattern",
				Description: "TODO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "TODO.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "TODO.",
				Type:        proto.ColumnType_INT,
			},

			// json fields
			{
				Name:        "groups",
				Description: "TODO.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "users",
				Description: "TODO.",
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
	client := connect(ctx, d)

	opts := &bitbucket.BranchRestrictionsOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	response, err := client.Repositories.BranchRestrictions.Gets(opts)
	if err != nil {
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
	}

	return nil, nil
}

func tableBitbucketBranchRestrictionGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketBranchRestrictionGet")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	restriction_id := d.KeyColumnQuals["id"].GetInt64Value()
	if repoFullName == "" {
		return nil, nil
	}
	owner, repoName := parseRepoFullName(repoFullName)

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repository_full_name should be in the format \"{workspace_slug}/{repo_slug}\"")
	}

	if restriction_id == 0 {
		return nil, nil
	}
	client := connect(ctx, d)

	opts := &bitbucket.BranchRestrictionsOptions{
		Owner:    owner,
		RepoSlug: repoName,
		ID:       strconv.Itoa(int(restriction_id)),
	}

	response, err := client.Repositories.BranchRestrictions.Get(opts)
	plugin.Logger(ctx).Error("tableBitbucketBranchRestrictionGet", "Responseeeeeeeeeee", response)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketBranchRestrictionGet", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	branchRestriction := new(bitbucket.BranchRestrictions)
	err = decodeJson(response, branchRestriction)
	plugin.Logger(ctx).Error("tableBitbucketBranchRestrictionGet", "After decodeeeeeeeeeeeeeee", branchRestriction)
	if err != nil {
		return nil, err
	}

	return branchRestriction, nil
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
	// Links           interface{} `json:"links,omitempty"`
	Value           *int        `json:"value,omitempty"`
	BranchMatchKind string      `json:"branch_match_kind,omitempty"`
	Groups          interface{} `json:"groups,omitempty"`
	BranchType      string      `json:"branch_type,omitempty"`
	Type            string      `json:"type,omitempty"`
}
