---
title: "Steampipe Table: bitbucket_project - Query Bitbucket Projects using SQL"
description: "Allows users to query Bitbucket Projects, specifically retrieving details about each project, including its ID, name, description, and more."
---

# Table: bitbucket_project - Query Bitbucket Projects using SQL

Bitbucket Projects is a feature within Bitbucket that allows users to group repositories into projects, making it easier to manage permissions and collaborate across multiple repositories. Projects provide a way to group repositories that are related to each other, such as those that are part of the same product or application. Projects also allow you to manage access control, so you can specify who has access to which repositories.

## Table Usage Guide

The `bitbucket_project` table provides insights into projects within Bitbucket. As a DevOps engineer or a software developer, explore project-specific details through this table, including project permissions, associated repositories, and other metadata. Utilize it to uncover information about projects, such as those with specific access controls, the associated repositories within each project, and the overall structure of your Bitbucket projects.

**Important Notes**
- You must specify the `workspace_slug` in the `where` or `join` clause (`where workspace_slug =`, `join bitbucket_project on workspace_slug =`) to query this table.

## Examples

### Get information about a specific project
Explore the details of a specific project by identifying key attributes such as the project name, unique identifier, workspace slug, owner's display name, and its privacy status. This can be particularly useful for project management and auditing purposes, allowing for a better understanding of project ownership and privacy settings.

```sql+postgres
select
  name,
  uuid,
  key as project_key,
  workspace_slug,
  owner_display_name,
  is_private
from
  bitbucket_project
where
  workspace_slug = 'np1981';
```

```sql+sqlite
select
  name,
  uuid,
  key as project_key,
  workspace_slug,
  owner_display_name,
  is_private
from
  bitbucket_project
where
  workspace_slug = 'np1981';
```