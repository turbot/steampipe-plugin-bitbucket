---
title: "Steampipe Table: bitbucket_commit - Query Bitbucket Commits using SQL"
description: "Allows users to query Bitbucket Commits, providing detailed information about each commit made in the Bitbucket repositories."
---

# Table: bitbucket_commit - Query Bitbucket Commits using SQL

Bitbucket is a web-based version control repository hosting service owned by Atlassian, for source code and development projects that use either Mercurial or Git revision control systems. Bitbucket offers both commercial plans and free accounts. It provides a way to manage and maintain versioning of source code, manage projects, work on your applications, and deploy them in a team environment.

## Table Usage Guide

The `bitbucket_commit` table provides insights into each commit made in the Bitbucket repositories. As a DevOps engineer or a developer, explore commit-specific details through this table, including commit messages, author details, and associated metadata. Utilize it to track changes, understand version history, and manage your development workflow more effectively.

## Examples

### List the commits in a repository
Discover the segments that have made changes in a specific repository. This can be used to track changes, understand the context of modifications, and identify the contributors involved.

```sql
select
  repository_full_name,
  hash,
  message,
  author_display_name,
  author_uuid
from
  bitbucket_commit
where
  repository_full_name = 'sayan97tb/stmp-rep';
```

### List the commits by a specific author
Explore the specific contributions made by an individual author within a particular repository. This could be useful for assessing their productivity or understanding the nature of their contributions.

```sql
select
  repository_full_name,
  hash,
  message,
  author_display_name,
  author_uuid
from
  bitbucket_commit
where
  repository_full_name = 'sayan97tb/stmp-rep'
  and author_display_name = 'sayan';
```