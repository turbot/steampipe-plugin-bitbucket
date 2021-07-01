# Table: bitbucket_project

Projects are used by teams to organize repositories.

The `bitbucket_project` table can be used to query information about ANY repository, and **you must specify which repository** in the where or join clause (`where workspace_slug =`, `join bitbucket_project on workspace_slug =`).

To list all the **your** projects use the `bitbucket_my_project` table instead. The `bitbucket_my_project` table will list projects you own, you collaborate on, or that belong to your workspaces.

## Examples

### Get information about a specific project

```sql
select
  name,
  uuid,
  key as project_key,
  workspace_slug,
  owner_display_name,
  is_private
from
  bitbucket_project
where
  workspace_slug = 'np1981';
```
