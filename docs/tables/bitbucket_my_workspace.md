# Table: bitbucket_my_workspace

A workspace is where you will create repositories, collaborate on your code, and organize different streams of work in your Bitbucket Cloud account. At this time, you'll be provided with one workspace and one workspace ID.

The `bitbucket_my_workspace` table will list the workspaces **that you are a member of**.

## Examples

### Basic info for the GitHub Workspaces to which you belong

```sql
select
  name as workspace,
  slug,
  uuid,
  is_private
from
  bitbucket_my_workspace;
```

### Get members details for workspaces you belong

```sql
select
  u.display_name as member_name,
  u.uuid as user_uuid,
  w.name as workspace,
  u.workspace_slug,
  u.account_id
from
  bitbucket_workspace_member as u,
  bitbucket_my_workspace as w
where
  w.slug = u.workspace_slug
order by
  w.slug;
```
