# Table: bitbucket_workspace_member

The `bitbucket_workspace_member` table list all members of Workspace. The table will list members only of the workspace that **you have required access to list members**.

**Note:** This table requires `workspace_slug` as the required input.

## Examples

### List members in the workspace

```sql
select
  display_name,
  uuid,
  account_id
from
  bitbucket_workspace_member
where
  workspace_slug = 'LalitFort';
```
