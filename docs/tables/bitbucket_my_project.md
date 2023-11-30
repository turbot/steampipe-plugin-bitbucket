---
title: "Steampipe Table: bitbucket_my_project - Query Bitbucket Projects using SQL"
description: "Allows users to query Bitbucket Projects, specifically providing details on project key, name, description, public access status and more."
---

# Table: bitbucket_my_project - Query Bitbucket Projects using SQL

Bitbucket Projects is a resource within Atlassian's Bitbucket that allows you to group your repositories into projects, making it easier to manage permissions and collaborate with your team. It provides a centralized way to manage and organize your repositories, giving you a higher level of control over your codebase. Bitbucket Projects helps you maintain a clean workspace and manage access to your repositories more efficiently.

## Table Usage Guide 

The `bitbucket_my_project` table provides insights into Bitbucket Projects within Atlassian's Bitbucket. As a software engineer or project manager, explore project-specific details through this table, including project keys, names, descriptions, and public access statuses. Utilize it to manage and organize your repositories, granting you a higher level of control over your codebase.

## Examples

### Get information about my projects
Explore your Bitbucket projects to gain insights into project names, unique identifiers, workspace slugs, ownership details, privacy status, and creation dates. This information can be useful for project management, tracking project ownership, and understanding the distribution of private and public projects.

```sql
select
  name,
  uuid,
  key as project_key,
  workspace_slug,
  owner_display_name,
  is_private,
  created
from
  bitbucket_my_project;
```

### List count of repositories by project
Analyze the distribution of repositories across various projects. This query can be used to understand the extent of codebase diversification within individual projects and identify those managed by specific owners.

```sql
select
  count(*),
  project_key,
  project_name,
  owner_display_name
from
  bitbucket_my_repository
group by
  project_key,
  project_name,
  owner_display_name
order by
  project_name;
```

### List count of my repositories by project
Explore the distribution of your repositories across different projects. This can help in understanding how your projects are structured and where the majority of your repositories are concentrated.

```sql
select
  count(r.project_key),
  p.key as project_key,
  p.name as project_name,
  p.owner_display_name
from
  bitbucket_my_project as p
  left join bitbucket_my_repository as r on r.project_key = p.key
group by
  p.key,
  p.name,
  p.owner_display_name
order by
  count,
  p.name;
```