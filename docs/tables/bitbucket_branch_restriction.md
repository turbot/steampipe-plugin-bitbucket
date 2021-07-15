# Table: bitbucket_branch_restriction

Bitbucket branch restrictions allow you to control the actions users can perform on a single branch, branch type, or branch pattern within a repository. Branch restrictions provide another level of security within Bitbucket Server, along with user authentication and project, repository and global permissions, that together allow you to control, or enforce, your own workflow or process.

The `bitbucket_branch_restriction` table can be used to query issues belonging to a repository, and **you must specify which repository** with `where repository_full_name='workspace/repository'`. Also **you must have admin access to the repository** in order to access the branch restrictions.

## Examples

### List the branch restrictions in a repository

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