---
title: "Steampipe Table: bitbucket_repository_src - Query Bitbucket repository files using SQL"
description: "Allows users to list the files and directories of a Bitbucket repository at a commit, and read the content of a single file."
---

# Table: bitbucket_repository_src - Query Bitbucket repository files using SQL

Bitbucket's source endpoint serves a repository's tree at any commit, branch or tag: a directory comes back as a listing, a file as its raw content.

## Table Usage Guide

The `bitbucket_repository_src` table lists the files in a repository and reads individual files without cloning it. As a security engineer, use it to find dependency manifests and lockfiles across repositories, or to check that a required file (a `CODEOWNERS`, a pipeline definition) exists and says what it should.

**Important Notes**
- You must specify `repository_full_name` in the `where` or `join` clause to query this table.
- Give `path` to read one file and its `content`, or leave it out to list a directory (`scope_path`, the repository root by default). The two cannot be combined.
- `commit` takes a commit hash, branch or tag and defaults to `HEAD`, the head of the main branch. Prefer a full hash: branch names containing `/` are ambiguous in this endpoint.
- `max_depth` lists nested directories in the same request (Bitbucket's default is 1), and `q` filters the listing server-side using Bitbucket's [filter syntax](https://developer.atlassian.com/cloud/bitbucket/rest/intro/#filtering), for example `type="commit_file" AND path~"package.json"`.
- Each page of up to 100 entries, and each file read, is one API request against the account's hourly limit.
- `content` is null for binary and non-UTF-8 files; `size` is still set.
- A path or commit that does not exist returns no rows. A token without read access to the repository gets `403`, which this table reports as an error.

## Examples

### List the top level of a repository
List the files and directories at the root of a repository's main branch.


```sql+postgres
select
  path,
  type,
  size
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo';
```

```sql+sqlite
select
  path,
  type,
  size
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo';
```

### Find dependency lockfiles anywhere in the tree
Find dependency manifests and lockfiles across the whole tree in one filtered listing.


```sql+postgres
select
  path,
  size,
  commit_hash
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo'
  and max_depth = 20
  and q = 'type="commit_file" AND (path~"package-lock.json" OR path~"yarn.lock" OR path~"go.mod")';
```

```sql+sqlite
select
  path,
  size,
  commit_hash
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo'
  and max_depth = 20
  and q = 'type="commit_file" AND (path~"package-lock.json" OR path~"yarn.lock" OR path~"go.mod")';
```

### Read one file at a commit
Read the content of a single file at a fixed commit.


```sql+postgres
select
  size,
  content
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo'
  and commit = '7b8b171c0dbfa8af665c24a63d91f2b185b8701b'
  and path = 'package.json';
```

```sql+sqlite
select
  size,
  content
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo'
  and commit = '7b8b171c0dbfa8af665c24a63d91f2b185b8701b'
  and path = 'package.json';
```

### List one directory on a branch
List the contents of one directory on a branch other than the main branch.


```sql+postgres
select
  path,
  type
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo'
  and commit = 'develop'
  and scope_path = 'services/api';
```

```sql+sqlite
select
  path,
  type
from
  bitbucket_repository_src
where
  repository_full_name = 'my-workspace/my-repo'
  and commit = 'develop'
  and scope_path = 'services/api';
```
