# Table: bitbucket_issue

Bitbucket issues are used to track ideas, enhancements, tasks, or bugs for work on Bitbucket.

The `bitbucket_issue` table can be used to query issues belonging to a repository, and **you must specify which repository** with `where repository_full_name='workspace/repository'`.

## Examples

### List the issues in a repository

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket';
```

### List the unassigned open issues in a repository

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
  and assignee_uuid is null
  and state = 'new';
```

### List the open issues in a repository assigned to a specific user

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
  and assignee_display_name = 'Lalit Bhardwaj'
  and state in ('new', 'open');
```

### Report of the number issues in a repository by author

```sql
select
  assignee_display_name,
  assignee_uuid,
  count(*) as num_issues
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
group by
  assignee_uuid,
  assignee_display_name
order by
  num_issues desc;
```
