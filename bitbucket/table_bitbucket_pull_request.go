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

func tableBitbucketPullRequest(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "bitbucket_pull_request",
		Description: "Issues are used to track ideas, enhancements, tasks, or bugs for work on Bitbucket.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableBitbucketPullRequestList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    tableBitbucketPullRequestGet,
		},
		Columns: []*plugin.Column{
			// top fields
			{
				Name:        "id",
				Description: "The pull request's immutable id.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "repository_full_name",
				Description: "The repository's full name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Destination.repository.full_name"),
			},
			{
				Name:        "created",
				Description: "Timestamp when pull request was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "state",
				Description: "A current state of the pull request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: "The title of pull request.",
				Type:        proto.ColumnType_STRING,
			},

			// other fields
			{
				Name:        "author_display_name",
				Description: "Display name of the author of this pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.display_name"),
			},
			{
				Name:        "author_uuid",
				Description: "UUID of the author of this pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Author.uuid"),
			},
			{
				Name:        "closed_by_display_name",
				Description: "Display name of the user who closed this pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClosedBy.display_name"),
			},
			{
				Name:        "closed_by_uuid",
				Description: "UUID of the user who closed of this pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClosedBy.uuid"),
			},
			{
				Name:        "edited",
				Description: "Timestamp when project was last edited.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "close_source_branch",
				Description: "",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "comment_count",
				Description: "The priority of the issue. Can be one of \"trivial\", \"minor\", \"major\", \"critical\", and \"blocker\".",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "Display name of the user issue is reported.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description"),
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
				Name:        "reason",
				Description: "No of the watchers on the issue.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "merge_commit",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "summary",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

func tableBitbucketPullRequestList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketPullRequestList")
	repoFullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)
	client := connect(ctx, d)

	opts := &bitbucket.PullRequestsOptions{
		Owner:    owner,
		RepoSlug: repoName,
	}

	response, err := client.Repositories.PullRequests.Gets(opts)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	pullRequestList := new(PullRequestList)
	err = decodeJson(response, pullRequestList)
	if err != nil {
		return nil, err
	}

	for _, issue := range pullRequestList.PullRequests {
		d.StreamListItem(ctx, issue)
	}

	return nil, nil
}

func tableBitbucketPullRequestGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketPullRequestGet")
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

type PullRequestList struct {
	ListResponse
	PullRequests []PullRequest `json:"values,omitempty"`
}

type PullRequest struct {
	Author            map[string]interface{} `json:"author,omitempty"`
	CloseSourceBranch bool                   `json:"close_source_branch,omitempty"`
	ClosedBy          map[string]interface{} `json:"closed_by,omitempty"`
	CommentCount      int                    `json:"comment_count,omitempty"`
	Created           *time.Time             `json:"created_on,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Destination       map[string]interface{} `json:"destination,omitempty"`
	ID                int                    `json:"id,omitempty"`
	Links             map[string]interface{} `json:"links,omitempty"`
	MergeCommit       map[string]interface{} `json:"merge_commit,omitempty"`
	Reason            string                 `json:"reason,omitempty"`
	Source            map[string]interface{} `json:"source,omitempty"`
	State             string                 `json:"state,omitempty"`
	Summary           map[string]interface{} `json:"summary,omitempty"`
	TaskCount         int                    `json:"task_count,omitempty"`
	Title             string                 `json:"title,omitempty"`
	Type              string                 `json:"type,omitempty"`
	Updated           *time.Time             `json:"updated_on,omitempty"`
}
