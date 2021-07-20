# Table: bitbucket_tag

Tags mark a specific commit at a point in your repository history. When you tag a commit, you're including all the changes before it. You can later compare tags to see the difference between two points in history. Tags are commonly used to mark release versions, with the release name as the tag name.

The `bitbucket_tag` table can be used to query information about ANY repository, and **you must specify which repository** in the where.

## Examples

### Get commits by tags

```sql
select
  name,
  repository_full_name,
  split_part(target -> 'links' ->> 'self', '/commit/', 2) as commit_hash
from
  bitbucket_tag
where
  repository_full_name = 'my-workspace/my-repo';
```
