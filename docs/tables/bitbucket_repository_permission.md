---
title: "Steampipe Table: bitbucket_repository_permission - Query Bitbucket repository permissions using SQL"
description: "Allows users to query the explicit user permissions held on the repositories of a Bitbucket workspace."
---

# Table: bitbucket_repository_permission - Query Bitbucket repository permissions using SQL

Bitbucket grants each user an explicit permission — `admin`, `write` or `read` — on a repository. This is the grant a least-privilege review reads: workspace membership alone says who is in the workspace, not what they can do to any particular repository.

## Table Usage Guide

The `bitbucket_repository_permission` table provides insight into who can do what to each repository in a workspace. As a security or compliance engineer, use it to find repositories with more administrators than they need, or users holding write access to repositories they do not work on.

**Important Notes**
- You must specify either `workspace_slug` or `repository_full_name` in the `where` or `join` clause to query this table.
- Bitbucket exposes these permissions per workspace, so a `repository_full_name` qualifier narrows that workspace listing.

## Examples

### List every explicit permission in a workspace

```sql+postgres
select
  repository_full_name,
  user_display_name,
  permission
from
  bitbucket_repository_permission
where
  workspace_slug = 'my-workspace';
```

```sql+sqlite
select
  repository_full_name,
  user_display_name,
  permission
from
  bitbucket_repository_permission
where
  workspace_slug = 'my-workspace';
```

### Find repository administrators

```sql+postgres
select
  repository_full_name,
  user_display_name,
  user_account_id
from
  bitbucket_repository_permission
where
  workspace_slug = 'my-workspace'
  and permission = 'admin';
```

```sql+sqlite
select
  repository_full_name,
  user_display_name,
  user_account_id
from
  bitbucket_repository_permission
where
  workspace_slug = 'my-workspace'
  and permission = 'admin';
```

### Count administrators per repository

```sql+postgres
select
  repository_full_name,
  count(*) as admin_count
from
  bitbucket_repository_permission
where
  workspace_slug = 'my-workspace'
  and permission = 'admin'
group by
  repository_full_name
order by
  admin_count desc;
```

```sql+sqlite
select
  repository_full_name,
  count(*) as admin_count
from
  bitbucket_repository_permission
where
  workspace_slug = 'my-workspace'
  and permission = 'admin'
group by
  repository_full_name
order by
  admin_count desc;
```

### List the permissions on one repository

```sql+postgres
select
  user_display_name,
  permission
from
  bitbucket_repository_permission
where
  repository_full_name = 'my-workspace/my-repo';
```

```sql+sqlite
select
  user_display_name,
  permission
from
  bitbucket_repository_permission
where
  repository_full_name = 'my-workspace/my-repo';
```
