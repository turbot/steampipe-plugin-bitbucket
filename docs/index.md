---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/bitbucket.svg"
brand_color: "#2684FF"
display_name: "Bitbucket"
short_name: "bitbucket"
description: "Steampipe plugin for querying repositories, issues, pull requests and more from Bitbucket."
og_description: "Query Bitbucket with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/bitbucket-social-graphic.png"
---

# Bitbucket + Steampipe

[Bitbucket](https://bitbucket.org) is a Git-based source code repository hosting service owned by Atlassian.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  name,
  uuid,
  is_private,
  full_name
from
  bitbucket_my_repository;
```

```
+----------------------------+----------------------------------------+------------+--------------------------------------+
| name                       | uuid                                   | is_private | full_name                            |
+----------------------------+----------------------------------------+------------+--------------------------------------+
| steampipe-plugin-bitbucket | {71624341-8873-4128-a356-f48c57c917e0} | true       | LalitFort/steampipe-plugin-bitbucket |
| repo2                      | {00ce5566-eba1-4a15-986d-85cc40f6b835} | true       | LalitFort/repo2                      |
+----------------------------+----------------------------------------+------------+--------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/bitbucket/tables)**

## Get started

### Install

Download and install the latest Bitbucket plugin:

```bash
steampipe plugin install bitbucket
```

### Credentials

| Item        | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| :---------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Bitbucket requires an [app password](https://bitbucket.org/account/settings/app-passwords/).                                                                                                                                                                                                                                                                                                                                                                                                         |
| Permissions | You must create app password with the following [scopes](https://developer.atlassian.com/cloud/bitbucket/rest/intro/#scopes):<br />&nbsp;&nbsp;&nbsp;&nbsp;- `account:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `issue:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `pullrequest:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `repository:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `snippet:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `webhook:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `wiki:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `workspace:read` |

### Configuration

Installing the latest bitbucket plugin will create a config file (`~/.steampipe/config/bitbucket.spc`) with a single connection named `bitbucket`:

```hcl
connection "bitbucket" {
  plugin = "bitbucket"

  # Bitbucket username.
  # Can also be set with the BITBUCKET_USERNAME environment variable.
  # username = "MyUsername"

  # Bitbucket app password, which can be created at https://bitbucket.org/account/settings/app-passwords/.
  # Can also be set with the BITBUCKET_PASSWORD environment variable.
  # password = "blHdmvlkFakeToken"

  # Base URL of your Bitbucket Server.
  # Defaults to "https://api.bitbucket.org/2.0".
  # Can also be set with the BITBUCKET_API_BASE_URL environment variable.
  # base_url = "https://api.bitbucket.org/2.0"
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-bitbucket
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
