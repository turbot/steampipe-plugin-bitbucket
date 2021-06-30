---
organization: Turbot
category: ["public cloud"]
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
| Credentials | Bitbucket requires an [app password](https://bitbucket.org/account/settings/app-passwords/), bitbucket instance baseurl (If private setup) and username for all requests.                                                                                                                                                                                                                                                                                                                                                                                                          |
| Permissions | You must create app password with the following [scopes](https://developer.atlassian.com/bitbucket/api/2/reference/meta/authentication#scopes-bbc):<br />&nbsp;&nbsp;&nbsp;&nbsp;- `workspace:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `account:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `snippet:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `wiki:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `repository:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `user:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `pullrequest:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `issue:read`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `webhook:read` |

### Configuration

Installing the latest bitbucket plugin will create a config file (`~/.steampipe/config/bitbucket.spc`) with a single connection named `bitbucket`:

```hcl
connection "bitbucket" {
  plugin   = "bitbucket"
  username = "LalitFort"
  password = "wOABk1jLlKktmtg43ZHNh9D12"
}
```

- `base_url` - The url of your bitbucket private instance.
- `username` - Bitbucket username.
- `password` - [App password](https://bitbucket.org/account/settings/app-passwords/) for bitbucket account.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-bitbucket
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
