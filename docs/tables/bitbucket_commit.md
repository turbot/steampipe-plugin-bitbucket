# Table: bitbucket_commit

Bitbucket commit is an operation which sends the latest changes of the source code to the Bitbucket repository, making these changes part of the head revision of the repository.

The `bitbucket_commit` table can be used to query commits belonging to a repository, and **you must specify which repository** with `where repository_full_name='workspace/repository'`.

## Examples

### List the commits in a repository

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
