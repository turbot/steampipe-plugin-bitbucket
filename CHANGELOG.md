## v0.6.0 [2023-08-04]

_Enhancements_

- Added the `default_reviewers` column to `bitbucket_repository` and `bitbucket_my_repository` tables. ([#64](https://github.com/turbot/steampipe-plugin-bitbucket/pull/64))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#72](https://github.com/turbot/steampipe-plugin-bitbucket/pull/72))

## v0.5.0 [2023-05-09]

_Enhancements_

- Added columns `account_status`, `created_on`, `is_staff` and `nickname` to `bitbucket_workspace_member` table. ([#71](https://github.com/turbot/steampipe-plugin-bitbucket/pull/71)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v0.4.0 [2023-04-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#69](https://github.com/turbot/steampipe-plugin-bitbucket/pull/69))

## v0.3.1 [2022-11-08]

_Bug fixes_

- Fixed the `docs/index.md` and the `config/bitbucket.spc` files to include the correct permission scope to query the tables. ([#66](https://github.com/turbot/steampipe-plugin-bitbucket/pull/66))

## v0.3.0 [2022-09-27]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#60](https://github.com/turbot/steampipe-plugin-bitbucket/pull/60))
- Recompiled plugin with Go version `1.19`. ([#60](https://github.com/turbot/steampipe-plugin-bitbucket/pull/60))

## v0.2.2 [2022-07-18]

_Bug fixes_

- Added pagination support to `bitbucket_branch` table to allow more than 10 results to be returned. ([#58](https://github.com/turbot/steampipe-plugin-bitbucket/pull/58))

## v0.2.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#53](https://github.com/turbot/steampipe-plugin-bitbucket/pull/53))

## v0.2.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#50](https://github.com/turbot/steampipe-plugin-bitbucket/pull/50))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#51](https://github.com/turbot/steampipe-plugin-bitbucket/pull/51))

## v0.1.0 [2021-11-23]

_Enhancements_

- Recompiled plugin with Go version 1.17 ([#47](https://github.com/turbot/steampipe-plugin-bitbucket/pull/47))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#46](https://github.com/turbot/steampipe-plugin-bitbucket/pull/46))

## v0.0.3 [2021-09-22]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v161--2021-09-21) ([#43](https://github.com/turbot/steampipe-plugin-bitbucket/pull/43))

## v0.0.2 [2021-07-22]

_What's new?_

- New tables added
  - [bitbucket_branch](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_branch) ([#36](https://github.com/turbot/steampipe-plugin-bitbucket/pull/36))
  - [bitbucket_branch_restriction](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_branch_restriction) ([#39](https://github.com/turbot/steampipe-plugin-bitbucket/pull/39))
  - [bitbucket_commit](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_commit) ([#26](https://github.com/turbot/steampipe-plugin-bitbucket/pull/26))
  - [bitbucket_tag](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_tag) ([#25](https://github.com/turbot/steampipe-plugin-bitbucket/pull/25))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v131--2021-07-15)

## v0.0.1 [2021-07-01]

_What's new?_

- New tables added

  - [bitbucket_issue](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_issue)
  - [bitbucket_my_project](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_my_project)
  - [bitbucket_my_repository](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_my_repository)
  - [bitbucket_my_workspace](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_my_workspace)
  - [bitbucket_project](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_project)
  - [bitbucket_pull_request](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_pull_request)
  - [bitbucket_repository](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_repository)
  - [bitbucket_workspace](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_workspace)
  - [bitbucket_workspace_member](https://hub.steampipe.io/plugins/turbot/bitbucket/tables/bitbucket_workspace_member)
