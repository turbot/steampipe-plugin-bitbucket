# Table: bitbucket_repository

A repository contains all of your project's files and each file's revision history.

The `bitbucket_repository` table can be used to query information about ANY repository, and **you must specify which repository** in the where or join clause (`where full_name=`, `join bitbucket_repository on full_name=`).

To list all the **your** repositories use the `bitbucket_my_repository` table instead. The `bitbucket_my_repository` table will list repositories you own, you collaborate on, or that belong to your organizations.

## Examples

### Get information about a specific repository

```sql
select
  name,
  uuid,
  full_name,
  owner_display_name,
  description
from
  bitbucket_repository
where
  full_name = 'bitbucketpipelines/official-pipes'
```
