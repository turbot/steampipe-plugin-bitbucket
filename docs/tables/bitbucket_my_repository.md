---
title: "Steampipe Table: bitbucket_my_repository - Query Bitbucket Repositories using SQL"
description: "Allows users to query Bitbucket Repositories, providing details about each repository including its name, project key, size, and more."
---

# Table: bitbucket_my_repository - Query Bitbucket Repositories using SQL

Bitbucket is a web-based version control repository hosting service owned by Atlassian, for source code and development projects that use either Mercurial or Git revision control systems. Bitbucket offers both commercial plans and free accounts. It provides a way for developers to manage and maintain their code, with features such as pull requests, branching, and in-line commenting.

## Table Usage Guide

The `bitbucket_my_repository` table provides insights into repositories within Bitbucket. As a developer or project manager, explore repository-specific details through this table, including repository name, project key, size, and other metadata. Utilize it to uncover information about repositories, such as those with the most recent commits, the size of each repository, and the associated project key.

## Examples

### List of repositories that you or your workspace owns
Explore the repositories that you or your workspace own, to manage and organize your projects better. This allows you to assess the ownership of different repositories, enhancing your control and coordination over them.

```sql
select
  name,
  uuid,
  full_name,
  owner_display_name
from
  bitbucket_my_repository
order by
  full_name;
```

### List your public repositories
Explore which of your Bitbucket repositories are publicly accessible. This can be useful to ensure sensitive information is not inadvertently exposed.

```sql
select
  name,
  is_private,
  full_name,
  owner_display_name
from
  bitbucket_my_repository
where
  not is_private;
```

### List the unassigned open issues in your repositories
Discover the segments that have unresolved issues in your repositories that are yet to be assigned. This is useful for identifying potential bottlenecks in your workflow, allowing you to take remedial action and improve efficiency.

```sql
select
  i.repository_full_name,
  i.id,
  i.title,
  i.state,
  i.assignee_display_name,
  i.assignee_uuid
from
  bitbucket_issue as i,
  bitbucket_my_repository as r
where
  repository_full_name = r.full_name
  and i.assignee_uuid is null
  and i.state = 'new';
```

### List details of the default reviewers of your repositories
Gain insights into who are set as default reviewers for your repositories, helpful in understanding the review process and ensuring the right people are involved in code reviews.

```sql
with default_reviewers as (
  select
  full_name as repository_name,
  r ->> 'AccountId' as reviewer_account_id,
  r ->> 'Uuid' as reviewer_uuid,
  r ->> 'DisplayName' as reviewer_display_name,
  r ->> 'Type' as reviewer_type
from
  bitbucket_my_repository,
  jsonb_array_elements(default_reviewers) as r
)
select
  repository_name,
  reviewer_account_id,
  reviewer_uuid,
  reviewer_display_name,
  reviewer_type
from
  default_reviewers;
```

### List the repositories without default reviewers
Discover the segments that have repositories without designated default reviewers. This is useful for identifying potential areas of oversight, ensuring that all repositories have appropriate review processes in place.

```sql
select
  name,
  uuid,
  full_name,
  owner_display_name
from
  bitbucket_my_repository
where
  default_reviewers is null;
```