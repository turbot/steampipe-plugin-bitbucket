---
title: "Steampipe Table: bitbucket_workspace - Query Bitbucket Workspaces using SQL"
description: "Allows users to query Bitbucket Workspaces, specifically the workspace details including the workspace ID, name, type, and associated metadata."
---

# Table: bitbucket_workspace - Query Bitbucket Workspaces using SQL

Bitbucket Workspaces are a key component of Bitbucket Cloud, providing a shared environment for teams to collaborate on code. Workspaces contain repositories, pull requests, and other resources, and are associated with a team or an individual user account. Each workspace has a unique ID, and can be customized with a name, type, and other metadata.

## Table Usage Guide

The `bitbucket_workspace` table provides insights into Bitbucket Workspaces within Bitbucket Cloud. As a DevOps engineer or a software developer, explore workspace-specific details through this table, including workspace ID, name, type, and associated metadata. Utilize it to manage and monitor workspaces, such as identifying workspaces with specific types, or retrieving workspace metadata for audit or reporting purposes.

**Important Notes**
- You must specify the `slug` in the `where` or `join` clause (`where slug=`, `join bitbucket_workspace on slug=`) to query this table.

## Examples

### Basic info for a workspace
Explore the key characteristics of a specific workspace, such as its name, unique identifier, and privacy status. This is useful for understanding the configuration and accessibility of your workspace within the Bitbucket platform.

```sql+postgres
select
  name as workspace,
  slug,
  uuid,
  is_private
from
  bitbucket_workspace
where
  slug = 'np1981';
```

```sql+sqlite
select
  name as workspace,
  slug,
  uuid,
  is_private
from
  bitbucket_workspace
where
  slug = 'np1981';
```