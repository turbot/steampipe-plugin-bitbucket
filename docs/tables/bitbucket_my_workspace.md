---
title: "Steampipe Table: bitbucket_my_workspace - Query Bitbucket Workspaces using SQL"
description: "Allows users to query Bitbucket Workspaces, specifically fetching information about the user's workspace, providing insights into the workspace's details, type, and associated resources."
---

# Table: bitbucket_my_workspace - Query Bitbucket Workspaces using SQL

Bitbucket Workspaces are a central location where teams can store and manage their code repositories. They provide a shared space for collaboration, code review, and version control. Workspaces hold repositories for the development team, providing a collaborative environment to manage and track source code changes.

## Table Usage Guide

The `bitbucket_my_workspace` table provides insights into the user's workspace within Bitbucket. As a DevOps engineer, explore workspace-specific details through this table, including the workspace's slug, type, and associated resources. Utilize it to uncover information about your workspace, such as its privacy settings, the owner's details, and the associated repositories.

## Examples

### Basic info for the Bitbucket workspaces to which you belong
Explore the Bitbucket workspaces you are a part of to understand their privacy settings. This is useful to ensure the correct level of access and security for your workspaces.

```sql
select
  name as workspace,
  slug,
  uuid,
  is_private
from
  bitbucket_my_workspace;
```

### Get members details for workspaces you belong
Explore the details of members within your shared workspaces. This can help you understand who else has access to the same resources, providing valuable context for collaboration and access management.
**Note:** Members will be listed for a workspace only if you have access to list them.


```sql
select
  u.display_name as member_name,
  u.uuid as user_uuid,
  w.name as workspace,
  u.workspace_slug,
  u.account_id
from
  bitbucket_workspace_member as u,
  bitbucket_my_workspace as w
where
  w.slug = u.workspace_slug
order by
  w.slug;
```