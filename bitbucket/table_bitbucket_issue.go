package bitbucket

import (
	"context"
	"time"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableBitbucketIssue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_issue",
		Description: "A Bitbucket project. Projects are used by teams to organize repositories.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketIssuesList,
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
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_STRING,
			},

			// other fields
			{
				Name:        "priority",
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kind",
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "votes",
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "watches",
				Description: "A usefule description for thr project.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "updated",
				Description: "Timestamp when project was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "edited",
				Description: "Timestamp when project was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "assignee_display_name",
				Description: "Display name of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Assignee.display_name"),
			},
			{
				Name:        "assignee_uuid",
				Description: "UUID of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Assignee.uuid"),
			},
			{
				Name:        "reporter_display_name",
				Description: "Display name of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Reporter.display_name"),
			},
			{
				Name:        "reporter_uuid",
				Description: "UUID of the owner of this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Reporter.uuid"),
			},
			{
				Name:        "self_link",
				Description: "A self link to this project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
		},
	}
}

func tableBitbucketIssuesList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketIssuesList")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)
	client := connect(ctx, d)

	// urlStr := client.GetApiBaseURL() + fmt.Sprintf("/workspaces/%s/projects", workspace.Slug)

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

type IssueList struct {
	Page     int     `json:"page,omitempty"`
	Pagelen  int     `json:"pagelen,omitempty"`
	MaxDepth int     `json:"maxDepth,omitempty"`
	Size     int     `json:"size,omitempty"`
	Next     string  `json:"next,omitempty"`
	Previous string  `json:"previous,omitempty"`
	Issues   []Issue `json:"values,omitempty"`
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
