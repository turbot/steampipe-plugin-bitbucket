---
title: "Steampipe Table: bitbucket_deploy_key - Query Bitbucket repository access keys using SQL"
description: "Allows users to query the access keys configured on a Bitbucket repository, including when each key was created and last used."
---

# Table: bitbucket_deploy_key - Query Bitbucket repository access keys using SQL

A Bitbucket access key (also called a deploy key) grants access to a single repository independently of any user account. Because a key is not tied to a person, it survives that person leaving, and it is not covered by workspace access reviews — which is what makes its age worth auditing.

## Table Usage Guide

The `bitbucket_deploy_key` table provides insight into the keys that can read a repository without signing in as a user. As a security or compliance engineer, use it to find keys that have never been used, or that are older than your rotation policy allows.

**Important Notes**
- You must specify the `repository_full_name` in the `where` or `join` clause (`where repository_full_name =`, `join bitbucket_deploy_key on repository_full_name =`) to query this table.
- Listing access keys requires the `admin:repository:bitbucket` scope. A token without it gets `403`, which this table reports as an error rather than as an empty list.

## Examples

### List the access keys on a repository

```sql+postgres
select
  id,
  label,
  created_on,
  last_used
from
  bitbucket_deploy_key
where
  repository_full_name = 'my-workspace/my-repo';
```

```sql+sqlite
select
  id,
  label,
  created_on,
  last_used
from
  bitbucket_deploy_key
where
  repository_full_name = 'my-workspace/my-repo';
```

### Find keys older than a year

```sql+postgres
select
  repository_full_name,
  label,
  created_on
from
  bitbucket_deploy_key
where
  repository_full_name = 'my-workspace/my-repo'
  and created_on < now() - interval '1 year';
```

```sql+sqlite
select
  repository_full_name,
  label,
  created_on
from
  bitbucket_deploy_key
where
  repository_full_name = 'my-workspace/my-repo'
  and created_on < datetime('now', '-1 year');
```

### Find keys that have never been used

```sql+postgres
select
  repository_full_name,
  label,
  created_on
from
  bitbucket_deploy_key
where
  repository_full_name = 'my-workspace/my-repo'
  and last_used is null;
```

```sql+sqlite
select
  repository_full_name,
  label,
  created_on
from
  bitbucket_deploy_key
where
  repository_full_name = 'my-workspace/my-repo'
  and last_used is null;
```

### Audit every key across a workspace

```sql+postgres
select
  k.repository_full_name,
  k.label,
  k.created_on,
  k.last_used
from
  bitbucket_repository as r
  join bitbucket_deploy_key as k on k.repository_full_name = r.full_name
where
  r.workspace_slug = 'my-workspace'
order by
  k.created_on;
```

```sql+sqlite
select
  k.repository_full_name,
  k.label,
  k.created_on,
  k.last_used
from
  bitbucket_repository as r
  join bitbucket_deploy_key as k on k.repository_full_name = r.full_name
where
  r.workspace_slug = 'my-workspace'
order by
  k.created_on;
```
