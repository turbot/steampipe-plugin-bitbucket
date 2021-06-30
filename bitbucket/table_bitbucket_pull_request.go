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
				Description: "The pull request's unique ID. Note that pull request IDs are only unique within their associated repository.",
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
				Description: "A current state of the pull request. Can be one of \"MERGED\", \"SUPERSEDED\", \"OPEN\" and \"DECLINED\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "title",
				Description: "Title of the pull request.",
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
				Name:        "branch_name",
				Description: "Name of the source branch for the pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Source.branch.name"),
			},
			{
				Name:        "comment_count",
				Description: "The number of comments for a specific pull request.",
				Type:        proto.ColumnType_INT,
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
				Name:        "close_source_branch",
				Description: "A boolean flag indicating if merging the pull request closes the source branch.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "description",
				Description: "Description of the pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "edited",
				Description: "Timestamp when pull request was last edited.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "merge_commit",
				Description: "Merge commit hash details for pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("MergeCommit.hash"),
			},
			{
				Name:        "summary",
				Description: "Summary details of the pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Summary.raw"),
			},
			{
				Name:        "self_link",
				Description: "A self link to this pull request.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Links.self.href"),
			},
			{
				Name:        "task_count",
				Description: "The number of open tasks for a specific pull request.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "type",
				Description: "Type of the Bitbucket resource. It will be always \"pullrequest\".",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "updated",
				Description: "Timestamp when pull request was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "participants",
				Description: "List of collaborators on the pull request.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     tableBitbucketPullRequestGet,
			},
			{
				Name:        "reviewers",
				Description: "List of reviewers of the pull request.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     tableBitbucketPullRequestGet,
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
		States:   []string{"merged", "open", "superseded", "declined"},
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

func tableBitbucketPullRequestGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableBitbucketPullRequestGet")

	var repoFullName string
	var issue_id int64

	if h.Item != nil {
		repoFullName = (h.Item.(PullRequest).Destination["repository"]).(map[string]interface{})["full_name"].(string)
		issue_id = int64(h.Item.(PullRequest).ID)
	} else {
		repoFullName = d.KeyColumnQuals["repository_full_name"].GetStringValue()
		issue_id = d.KeyColumnQuals["id"].GetInt64Value()
	}

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

	opts := &bitbucket.PullRequestsOptions{
		Owner:    owner,
		RepoSlug: repoName,
		ID:       strconv.Itoa(int(issue_id)),
	}

	response, err := client.Repositories.PullRequests.Get(opts)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		plugin.Logger(ctx).Error("tableBitbucketPullRequestGet", "Error", err)
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	pullRequest := new(PullRequest)
	err = decodeJson(response, pullRequest)
	if err != nil {
		return nil, err
	}

	return pullRequest, nil
}

type PullRequestList struct {
	ListResponse
	PullRequests []PullRequest `json:"values,omitempty"`
}

type PullRequest struct {
	Author            map[string]interface{}   `json:"author,omitempty"`
	CloseSourceBranch bool                     `json:"close_source_branch,omitempty"`
	ClosedBy          map[string]interface{}   `json:"closed_by,omitempty"`
	CommentCount      int                      `json:"comment_count,omitempty"`
	Created           *time.Time               `json:"created_on,omitempty"`
	Description       string                   `json:"description,omitempty"`
	Destination       map[string]interface{}   `json:"destination,omitempty"`
	ID                int                      `json:"id,omitempty"`
	Links             map[string]interface{}   `json:"links,omitempty"`
	MergeCommit       map[string]interface{}   `json:"merge_commit,omitempty"`
	Source            map[string]interface{}   `json:"source,omitempty"`
	Reviewers         []map[string]interface{} `json:"reviewers,omitempty"`
	Participants      []map[string]interface{} `json:"participants,omitempty"`
	State             string                   `json:"state,omitempty"`
	Summary           map[string]interface{}   `json:"summary,omitempty"`
	TaskCount         int                      `json:"task_count,omitempty"`
	Title             string                   `json:"title,omitempty"`
	Type              string                   `json:"type,omitempty"`
	Updated           *time.Time               `json:"updated_on,omitempty"`
}
