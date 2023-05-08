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
  i.repository_full_name,
  i.id,
  i.title,
  i.state,
  i.assignee_display_name,
  i.assignee_uuid
from
  bitbucket_issue as i,
  bitbucket_my_repository as r
where
  repository_full_name = r.full_name
  and i.assignee_uuid is null
  and i.state = 'new';
```

### List details of the default reviewers of your repositories

```sql
with default_reviewers as (
  select
  full_name as repository_name,
  r ->> 'AccountId' as reviewer_account_id,
  r ->> 'Uuid' as reviewer_uuid,
  r ->> 'DisplayName' as reviewer_display_name,
  r ->> 'Type' as reviewer_type
from
  bitbucket_my_repository,
  jsonb_array_elements(default_reviewers) as r
)
select
  repository_name,
  reviewer_account_id,
  reviewer_uuid,
  reviewer_display_name,
  reviewer_type
from
  default_reviewers;
```

### List the repositories without default reviewers

```sql
select
  name,
  uuid,
  full_name,
  owner_display_name
from
  bitbucket_my_repository
where
  default_reviewers is null;
```