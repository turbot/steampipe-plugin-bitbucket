---
title: "Steampipe Table: bitbucket_repository - Query Bitbucket Repositories using SQL"
description: "Allows users to query Bitbucket Repositories, specifically the repository name, UUID, project key, and other related details, providing insights into the repository's metadata and status."
---

# Table: bitbucket_repository - Query Bitbucket Repositories using SQL

Bitbucket is a web-based version control repository hosting service owned by Atlassian, for source code and development projects that use either Mercurial or Git revision control systems. Bitbucket offers both commercial plans and free accounts. It provides free private repositories for small teams and also offers features such as pull requests, branching, and in-line commenting.

## Table Usage Guide

The `bitbucket_repository` table provides insights into Bitbucket repositories within Atlassian's Bitbucket service. As a DevOps engineer, explore repository-specific details through this table, including repository name, UUID, project key, and other related details. Utilize it to uncover information about repositories, such as their status, metadata, and the projects they belong to.

**Important Notes**
- You must specify either `workspace_slug` or `full_name` in the `where` or `join` clause to query this table. `workspace_slug` lists a workspace's repositories; `full_name` fetches one.
- Use `workspace_slug` rather than `bitbucket_my_repository` to enumerate repositories. `bitbucket_my_repository` lists the authenticated account's workspaces first, and the endpoint it uses for that (`/2.0/workspaces`) has been withdrawn by Atlassian.

## Examples

### List the repositories in a workspace
Enumerate every repository a workspace holds, with the fields that tell a live repository from an abandoned one.

```sql+postgres
select
  full_name,
  is_private,
  fork_policy,
  updated_on,
  size
from
  bitbucket_repository
where
  workspace_slug = 'my-workspace';
```

```sql+sqlite
select
  full_name,
  is_private,
  fork_policy,
  updated_on,
  size
from
  bitbucket_repository
where
  workspace_slug = 'my-workspace';
```

### Find repositories with no activity in 180 days

```sql+postgres
select
  full_name,
  updated_on
from
  bitbucket_repository
where
  workspace_slug = 'my-workspace'
  and updated_on < now() - interval '180 days';
```

```sql+sqlite
select
  full_name,
  updated_on
from
  bitbucket_repository
where
  workspace_slug = 'my-workspace'
  and updated_on < datetime('now', '-180 days');
```


### Get information about a specific repository
Explore specific information about a designated repository to understand its owner, unique identifier, and description. This is useful in gaining insights into the repository's details without having to manually search through Bitbucket.

```sql+postgres
select
  name,
  uuid,
  full_name,
  owner_display_name,
  description
from
  bitbucket_repository
where
  full_name = 'bitbucketpipelines/official-pipes'
```

```sql+sqlite
select
  name,
  uuid,
  full_name,
  owner_display_name,
  description
from
  bitbucket_repository
where
  full_name = 'bitbucketpipelines/official-pipes'
```