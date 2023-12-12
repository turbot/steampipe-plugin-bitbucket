---
title: "Steampipe Table: bitbucket_workspace_member - Query Bitbucket Workspace Members using SQL"
description: "Allows users to query Bitbucket Workspace Members, providing insights into workspace member details, roles, and permissions."
---

# Table: bitbucket_workspace_member - Query Bitbucket Workspace Members using SQL

Bitbucket is a Git-based source code repository hosting service owned by Atlassian. The workspace member in Bitbucket is a user who has been granted access to a workspace. Workspace members can have different roles and permissions based on their level of access.

## Table Usage Guide

The `bitbucket_workspace_member` table provides insights into workspace members within Bitbucket. As a DevOps engineer, explore member-specific details through this table, including roles, permissions, and associated metadata. Utilize it to uncover information about members, such as their roles, the level of access, and the verification of permissions.

**Important Notes**
- You must specify the `workspace_slug` in the `where` or `join` clause (`where workspace_slug=`, `join bitbucket_workspace_member on workspace_slug=`) to query this table.

## Examples

### List members in the workspace
Explore which members are part of a specific workspace, providing an overview of who has access to shared resources and projects. This is particularly useful for workspace administrators to manage access and permissions effectively.

```sql+postgres
select
  display_name,
  uuid,
  account_id
from
  bitbucket_workspace_member
where
  workspace_slug = 'LalitFort';
```

```sql+sqlite
select
  display_name,
  uuid,
  account_id
from
  bitbucket_workspace_member
where
  workspace_slug = 'LalitFort';
```