---
title: "Steampipe Table: bitbucket_tag - Query Bitbucket Tags using SQL"
description: "Allows users to query Bitbucket Tags, specifically the name, repository, and other details of tags in a repository, providing insights into version control and codebase organization."
---

# Table: bitbucket_tag - Query Bitbucket Tags using SQL

Bitbucket Tags are reference points in your code, marking specific points in your repository's history. It's a way to capture a point in history that marks a significant event, such as a release. Bitbucket tags make it easier to navigate through the repository, identify current versions, and manage your codebase effectively. 

## Table Usage Guide

The `bitbucket_tag` table provides insights into tags within Bitbucket repositories. As a DevOps engineer, explore tag-specific details through this table, including the repository it belongs to, its name, and its associated details. Utilize it to manage and organize your codebase effectively, identify current versions of your application, and navigate through your repository with ease.

**Important Notes**
- You must specify the `repository_full_name` in the `where` clause to query this table.

## Examples

### Get commits by tags
Explore which commits are associated with tags in a specific repository. This can be helpful in understanding the relationship between commits and tags in your Bitbucket workspace, aiding in version control and project management.

```sql+postgres
select
  name,
  repository_full_name,
  split_part(target -> 'links' ->> 'self', '/commit/', 2) as commit_hash
from
  bitbucket_tag
where
  repository_full_name = 'my-workspace/my-repo';
```

```sql+sqlite
Error: SQLite does not support split_part function.
```