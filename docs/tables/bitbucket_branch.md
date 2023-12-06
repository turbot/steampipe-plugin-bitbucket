---
title: "Steampipe Table: bitbucket_branch - Query Bitbucket Branches using SQL"
description: "Allows users to query Bitbucket Branches, specifically providing details about each branch in a Bitbucket repository, including its name, type, and associated commit."
---

# Table: bitbucket_branch - Query Bitbucket Branches using SQL

Bitbucket is a Git-based source code repository hosting service owned by Atlassian. Bitbucket branches represent independent lines of development within a repository. They allow developers to work on features or fixes without interfering with the main (or 'master') codebase, and can be merged back into the main codebase when the work is complete.

## Table Usage Guide

The `bitbucket_branch` table provides insights into branches within Bitbucket. As a developer or DevOps engineer, explore branch-specific details through this table, including branch name, type, and associated commit. Utilize it to uncover information about branches, such as their current status, the associated commit, and the overall structure of development within a repository.

**Important Notes**

- You must specify the `repository_full_name` in the `where` clause in order to query this table.
- You must have admin access to the repository in order to access the branch restrictions.

## Examples

### List branches in a repository
Explore the branches within a specific repository to understand their default merge strategies and types. This can be useful for managing and optimizing your project's workflow within Bitbucket.

```sql+postgres
select
  name,
  repository_full_name,
  default_merge_strategy,
  type
from
  bitbucket_branch
where
  repository_full_name = 'souravthe/test1';
```

```sql+sqlite
select
  name,
  repository_full_name,
  default_merge_strategy,
  type
from
  bitbucket_branch
where
  repository_full_name = 'souravthe/test1';
```

### List branches with `merge_commit` default merge strategy
Explore which branches within a specific repository utilize the 'merge_commit' as their default merge strategy. This is useful for understanding the merging practices and strategies within your project.

```sql+postgres
select
  name,
  repository_full_name,
  default_merge_strategy,
  type
from
  bitbucket_branch
where
  repository_full_name = 'souravthe/test1'
  and default_merge_strategy = 'merge_commit';
```

```sql+sqlite
select
  name,
  repository_full_name,
  default_merge_strategy,
  type
from
  bitbucket_branch
where
  repository_full_name = 'souravthe/test1'
  and default_merge_strategy = 'merge_commit';
```