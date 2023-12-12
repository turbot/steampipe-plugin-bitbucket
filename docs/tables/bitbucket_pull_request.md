---
title: "Steampipe Table: bitbucket_pull_request - Query Bitbucket Pull Requests using SQL"
description: "Allows users to query Bitbucket Pull Requests, providing detailed information on specific pull request data such as ID, title, description, and status."
---

# Table: bitbucket_pull_request - Query Bitbucket Pull Requests using SQL

Bitbucket Pull Requests is a feature within Bitbucket that allows developers to review code and discuss changes before merging into the main code base. It helps in maintaining code quality, fostering collaboration, and promoting a more transparent development process. With pull requests, teams can ensure that only thoroughly reviewed and approved code gets merged, reducing the risk of introducing bugs.

## Table Usage Guide

The `bitbucket_pull_request` table provides insights into pull requests within Bitbucket. As a developer or team lead, you can explore pull request details through this table, including the author, reviewers, status, and associated metadata. Utilize it to monitor the progress of code reviews, track the status of pull requests, and ensure efficient collaboration within your team.

**Important Notes**
- You must specify the `repository_full_name` in the `where` clause to query this table.

## Examples

### List open pull requests in a repository
Explore which pull requests are currently open in a specific repository. This is useful to track ongoing changes and contributions to a project.

```sql+postgres
select
  repository_full_name,
  id,
  title,
  state,
  branch_name,
  author_display_name
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and state = 'OPEN';
```

```sql+sqlite
select
  repository_full_name,
  id,
  title,
  state,
  branch_name,
  author_display_name
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and state = 'OPEN';
```

### List the pull requests for a repository that have been merged in the last week
This query is useful for tracking recent changes in a specific repository, specifically identifying which pull requests have been merged in the past week. It provides valuable insights into the recent activity and contributor involvement, aiding in project management and review processes.

```sql+postgres
select
  repository_full_name,
  id,
  title,
  state,
  updated as merged_at,
  closed_by_display_name as merged_by_name,
  closed_by_uuid as merged_by_uuid
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and state = 'MERGED'
  and updated >= (current_date - interval '7' day)
order by
  updated desc;
```

```sql+sqlite
select
  repository_full_name,
  id,
  title,
  state,
  updated as merged_at,
  closed_by_display_name as merged_by_name,
  closed_by_uuid as merged_by_uuid
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and state = 'MERGED'
  and updated >= date('now','-7 day')
order by
  updated desc;
```

### List the open PRs in a repository assigned to a specific user
Determine the open project requests assigned to a specific user within a particular repository. This is useful in tracking the workload and contributions of individual team members.

```sql+postgres
select
  repository_full_name,
  id,
  title,
  state,
  author_display_name,
  author_uuid
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and author_display_name = 'Tom Gibson'
  and state = 'OPEN';
```

```sql+sqlite
select
  repository_full_name,
  id,
  title,
  state,
  author_display_name,
  author_uuid
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and author_display_name = 'Tom Gibson'
  and state = 'OPEN';
```