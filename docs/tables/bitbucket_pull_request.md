# Table: bitbucket_pull_request

Bitbucket pull requests let you tell others about changes you've pushed to a branch in a repository. Once a pull request is opened, you can discuss and review the potential changes with collaborators and add follow-up commits before your changes are merged into the base branch.

The `bitbucket_pull_request` table can be used to query pull requests belonging to a repository. **You must specify which repository** in a `where` or `join` clause (`where repository_full_name='`, `join bitbucket_pull_request on repository_full_name=`).

## Examples

### List open pull requests in a repository

```sql
select
  repository_full_name,
  id,
  title,
  state,
  branch_name,
  author_display_name,
from
  bitbucket_pull_request
where
  repository_full_name = 'bitbucketpipelines/official-pipes'
  and state = 'OPEN';
```

### List the pull requests for a repository that have been merged in the last week

```sql
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

### List the open PRs in a repository assigned to a specific user

```sql
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
