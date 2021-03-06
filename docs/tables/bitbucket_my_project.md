# Table: bitbucket_my_project

Projects are used by teams to organize repositories.

The `bitbucket_my_repository` table will lists projects your user have access or that belong to your workspaces.

To query **ANY** project, including public projects, use the `bitbucket_my_repository` table.

## Examples

### Get information about my projects

```sql
select
  name,
  uuid,
  key as project_key,
  workspace_slug,
  owner_display_name,
  is_private,
  created
from
  bitbucket_my_project;
```

### List count of repositories by project

```sql
select
  count(*),
  project_key,
  project_name,
  owner_display_name
from
  bitbucket_my_repository
group by
  project_key,
  project_name,
  owner_display_name
order by
  project_name;
```

### List count of my repositories by project

```sql
select
  count(r.project_key),
  p.key as project_key,
  p.name as project_name,
  p.owner_display_name
from
  bitbucket_my_project as p
  left join bitbucket_my_repository as r on r.project_key = p.key
group by
  p.key,
  p.name,
  p.owner_display_name
order by
  count,
  p.name;
```
