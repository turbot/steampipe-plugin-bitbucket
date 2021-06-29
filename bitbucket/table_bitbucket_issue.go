package bitbucket

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableBitbucketIssue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_issue",
		Description: "Issues are used to track ideas, enhancements, tasks, or bugs for work on Bitbucket.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketIssuesList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    tableBitbucketIssueGet,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "id",
				Description: "The issues's immutable id.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "repository_full_name",
				Description: "The project's immutable id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Repository.full_name"),
			},
			{
				Name:        "created",
				Description: "Timestamp when issue was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state",
				Description: "A current state of the issue. Can we one of new \"open\", \"resolved\",\"on hold\", \"invalid\", \"duplicate\", \"wontfix\" and \"closed\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: "The issue title.",
				Type:        proto.ColumnType_STRING,
			},

			// other fields
			{
				Name:        "assignee_display_name",
				Description: "Display name of the assignee of this issue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Assignee.display_name"),
			},
			{
				Name:        "assignee_uuid",
				Description: "UUID of the assignee of this issue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Assignee.uuid"),
			},
			{
				Name:        "edited",
				Description: "Timestamp when project was last edited.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "kind",
				Description: "The kind of the issue. Can be one of \"bug\", \"enhancement\", \"proposal\", and \"task\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "priority",
				Description: "The priority of the issue. Can be one of \"trivial\", \"minor\", \"major\", \"critical\", and \"blocker\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reporter_display_name",
				Description: "Display name of the user issue is reported.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Reporter.display_name"),
			},
			{
				Name:        "reporter_uuid",
				Description: "UUID of the user issue is reported.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Reporter.uuid"),
			},
			{
				Name:        "self_link",
				Description: "A self link to this issue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource. It will be always \"issue\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated",
				Description: "Timestamp when issue was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "votes",
				Description: "Number of the upvotes on the issue.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "watches",
				Description: "No of the watchers on the issue.",
				Type:        proto.ColumnType_INT,
			},

			// json fields
			{
				Name:        "component",
				Description: "Content object of the issue with the rendering type details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "content",
				Description: "Version is a point in project or product timeline.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "milestone",
				Description: "A milestone is a subset of a version. It is a point that a development team works towards.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "version",
				Description: "Version is a point in project or product timeline.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

func tableBitbucketIssuesList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketIssuesList")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)
	client := connect(ctx, d)

	opts := &bitbucket.IssuesOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	response, err := client.Repositories.Issues.Gets(opts)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	issueList := new(IssueList)
	err = decodeJson(response, issueList)
	if err != nil {
		return nil, err
	}

	for _, issue := range issueList.Issues {
		d.StreamListItem(ctx, issue)
	}

	return nil, nil
}

func tableBitbucketIssueGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketIssueGet")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	issue_id := d.KeyColumnQuals["id"].GetInt64Value()

	if repoFullName == "" {
		return nil, nil
	}
	owner, repoName := parseRepoFullName(repoFullName)

	if owner == "" || repoName == "" {
		return nil, fmt.Errorf("repository_full_name should be in the format \"{workspace_slug}/{repo_slug}\"")
	}

	if issue_id == 0 {
		return nil, nil
	}
	client := connect(ctx, d)

	opts := &bitbucket.IssuesOptions{
		Owner:    owner,
		RepoSlug: repoName,
		ID:       strconv.Itoa(int(issue_id)),
	}

	response, err := client.Repositories.Issues.Get(opts)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("getEpic", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	issue := new(Issue)
	err = decodeJson(response, issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

type IssueList struct {
	ListResponse
	Issues []Issue `json:"values,omitempty"`
}

type Issue struct {
	Assignee   map[string]interface{} `json:"assignee,omitempty"`
	Component  map[string]interface{} `json:"component,omitempty"`
	Content    map[string]interface{} `json:"content,omitempty"`
	Created    *time.Time             `json:"created_on,omitempty"`
	Edited     *time.Time             `json:"edited_on,omitempty"`
	ID         int                    `json:"id,omitempty"`
	Kind       string                 `json:"kind,omitempty"`
	Links      map[string]interface{} `json:"links,omitempty"`
	Milestone  map[string]interface{} `json:"milestone,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Priority   string                 `json:"priority,omitempty"`
	Reporter   map[string]interface{} `json:"reporter,omitempty"`
	Repository map[string]interface{} `json:"repository,omitempty"`
	State      string                 `json:"state,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Updated    *time.Time             `json:"updated_on,omitempty"`
	Version    map[string]interface{} `json:"version,omitempty"`
	Votes      int                    `json:"votes,omitempty"`
	Watches    int                    `json:"watches,omitempty"`
	Workspace  bitbucket.Workspace    `json:"workspace,omitempty"`
}
