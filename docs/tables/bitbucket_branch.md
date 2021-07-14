# Table: bitbucket_branch

A branch represents an independent line of development. Branches serve as an abstraction for the edit/stage/commit process. You can think of them as a way to request a brand new working directory, staging area, and project history. New commits are recorded in the history for the current branch, which results in a fork in the history of the project.

The `bitbucket_branch` table can be used to query branches belonging to a repository, and **you must specify which repository** with `where repository_full_name='workspace/repository'`.

## Examples

### List branches in a repository

```sql
select
  name,
  repository_full_name,
  default_merge_strategy,
  type
from
  bitbucket_branch
where
  repository_full_name = 'souravthe/test1';
```

### List branches where default merge strategy is merge_commit

```sql
select
  name,
  repository_full_name,
  default_merge_strategy,
  type
from
  bitbucket_branch
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
  and default_merge_strategy = 'merge_commit';
```
