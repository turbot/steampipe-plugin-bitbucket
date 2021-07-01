# Table: bitbucket_workspace

A workspace is where you will create repositories, collaborate on your code, and organize different streams of work in your Bitbucket Cloud account. At this time, you'll be provided with one workspace and one workspace ID.

The `bitbucket_workspace` table can be used to query information about ANY workspace, and **you must specify which workspace** in the where or join clause (`where slug=`, `join bitbucket_workspace on slug=`).

## Examples

### Basic info for a workspace

```sql
select
  name as workspace,
  slug,
  uuid,
  is_private
from
  bitbucket_workspace
where
  slug = 'np1981';
```
