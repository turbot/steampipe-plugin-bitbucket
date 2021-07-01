# Table: bitbucket_my_repository

A repository contains all of your project's files and each file's revision history.

You can own repositories individually, or you can share ownership of repositories with other people in an organization. The `bitbucket_my_repository` table will list repositories your username own or that belong to your workspaces.

To query **ANY** repository, including public repos, use the `bitbucket_repository` table.

## Examples

### List of repositories that you or your workspace owns

```sql
select
  name,
  uuid,
  full_name,
  owner_display_name
from
  bitbucket_my_repository
order by
  full_name;
```

### List your public repositories

```sql
select
  name,
  is_private,
  full_name,
  owner_display_name
from
  bitbucket_my_repository
where
  not is_private;
```

### List the unassigned open issues in your repositories

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue,
  bitbucket_my_repository as r
where
  repository_full_name = r.full_name
  and assignee_uuid is null
  and state = 'new';
```
