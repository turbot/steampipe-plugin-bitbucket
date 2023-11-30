---
title: "Steampipe Table: bitbucket_branch_restriction - Query Bitbucket Branch Restrictions using SQL"
description: "Allows users to query Bitbucket Branch Restrictions, specifically the restrictions applied to a repository's branches, providing insights into the level of access and operations permitted."
---

# Table: bitbucket_branch_restriction - Query Bitbucket Branch Restrictions using SQL

Bitbucket Branch Restrictions is a feature within Bitbucket that allows you to control what actions users can perform on a repository's branches. This includes restricting who can push to a branch, who can merge or delete a branch, and enforcing a minimum number of reviewers for pull requests. It provides a layer of security and control over your codebase, ensuring only authorized users can make significant changes.

## Table Usage Guide

The `bitbucket_branch_restriction` table provides insights into the restrictions applied to a repository's branches within Bitbucket. As a DevOps engineer or repository administrator, explore branch-specific details through this table, including who can perform certain operations and the restrictions in place. Utilize it to uncover information about branch restrictions, such as those with specific user access, the operations permitted on a branch, and the enforcement of review policies.

## Examples

### List the branch restrictions in a repository
Explore the restrictions placed on different branches within a specific repository. This is useful for understanding the limitations and rules that have been set, which can help in managing code changes and merges effectively.

```sql
select
  repository_full_name,
  id,
  self_link,
  kind,
  value,
  type
from
  bitbucket_branch_restriction
where
  repository_full_name = 'sayan97tb/stmp-rep';
```

### List the branch restrictions having pattern for branch names
Discover the segments that have specific naming patterns for branch restrictions within a given repository. This is useful to maintain naming conventions and manage access control in a systematic manner.

```sql
select
  repository_full_name,
  id,
  self_link,
  pattern
from
  bitbucket_branch_restriction
where
  repository_full_name = 'sayan97tb/stmp-rep'
  and pattern = 'test-*';
```

### List the branch restrictions associated with a specific user
Determine the specific branch restrictions linked to a particular user within a Bitbucket repository. This allows you to understand and manage access and modification rights for individual users, enhancing repository security and collaboration efficiency.

```sql
select
  repository_full_name,
  id,
  self_link,
  pattern,
  u ->> 'display_name' as user_name
from
  bitbucket_branch_restriction,
  jsonb_array_elements(users) as u
where
  repository_full_name = 'sayan97tb/stmp-rep'
  and u ->> 'display_name' = 'sayan';
```

### List the branch restrictions having 'branching_model' branch_match_kind
Explore the restrictions on specific branches within a repository to understand its access and modification permissions. This is useful for managing and securing your codebase by controlling who can make changes to specific branches.

```sql
select
  repository_full_name,
  id,
  self_link,
  pattern,
  branch_match_kind
from
  bitbucket_branch_restriction
where
  repository_full_name = 'sayan97tb/stmp-rep'
  and branch_match_kind = 'branching_model';
```