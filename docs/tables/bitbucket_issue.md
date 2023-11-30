---
title: "Steampipe Table: bitbucket_issue - Query Bitbucket Issues using SQL"
description: "Allows users to query Bitbucket Issues, specifically providing details about each issue such as its ID, title, type, priority, status, and assignee, offering insights into project management and task distribution."
---

# Table: bitbucket_issue - Query Bitbucket Issues using SQL

Bitbucket Issues is a feature within Bitbucket that allows you to track and manage tasks, enhancements, and bugs for your projects. It provides a centralized platform to manage issues for various Bitbucket repositories, including details such as issue ID, title, type, priority, status, and assignee. Bitbucket Issues helps you stay informed about the progress and distribution of tasks within your project.

## Table Usage Guide

The `bitbucket_issue` table provides insights into issue management within Bitbucket. As a project manager or developer, explore issue-specific details through this table, including issue type, priority, status, and assignee. Utilize it to uncover information about issues, such as their distribution among team members, the status of various tasks, and the prioritization of enhancements and bugs.

## Examples

### List the issues in a repository
Explore which issues are currently present in a specific project repository. This is useful for project managers who need to assess the status and assignment of issues for effective project management.

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue
where
  repository_full_name = 'gamesaucer/mono-ui';
```

### List the unassigned open issues in a repository
Explore which open issues in a specific repository have not been assigned yet, enabling efficient task allocation and workload management.

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
  and assignee_uuid is null
  and state in ('new','open');
```

### List the open issues in a repository assigned to a specific user
Explore which open issues in a specific repository are assigned to a particular user. This can be useful for project managers to track individual workloads and identify any potential bottlenecks in project progression.

```sql
select
  repository_full_name,
  id,
  title,
  state,
  assignee_display_name,
  assignee_uuid
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
  and assignee_display_name = 'Lalit Bhardwaj'
  and state in ('new', 'open');
```

### Report of the number issues in a repository by author
Analyze the distribution of issues in a specific repository based on the author to understand their individual contributions and identify any patterns or anomalies.

```sql
select
  assignee_display_name,
  assignee_uuid,
  count(*) as num_issues
from
  bitbucket_issue
where
  repository_full_name = 'LalitFort/steampipe-plugin-bitbucket'
group by
  assignee_uuid,
  assignee_display_name
order by
  num_issues desc;
```

### List the unassigned open issues in your repositories
Explore the open issues in your repositories that have not been assigned to anyone. This is useful in identifying areas that need attention or tasks that are yet to be allocated to team members.

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
  and assignee_uuid is null
  and state in ('new', 'open');
```