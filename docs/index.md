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
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# Bitbucket + Steampipe

[Bitbucket](https://bitbucket.org) is a Git-based source code repository hosting service owned by Atlassian.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

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
| Credentials | Bitbucket requires an [API token](https://support.atlassian.com/bitbucket-cloud/docs/create-an-api-token/). Set `username` to your **Atlassian account email** (listed under *Email Aliases* in your Bitbucket Personal settings), not your Bitbucket username, and `password` to the API token. |
| Permissions | Create the token via *Atlassian account settings → Security → Create and manage API tokens → Create API token with scopes*, select **Bitbucket** as the app, and grant these [scopes](https://support.atlassian.com/bitbucket-cloud/docs/api-token-permissions/):<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:user:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:workspace:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:project:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:repository:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:pullrequest:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:issue:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:snippet:bitbucket`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:webhook:bitbucket` |

### Configuration

Installing the latest bitbucket plugin will create a config file (`~/.steampipe/config/bitbucket.spc`) with a single connection named `bitbucket`:

```hcl
connection "bitbucket" {
  plugin = "bitbucket"

  # Atlassian account email (the address listed under "Email Aliases" in your Bitbucket Personal settings).
  # Can also be set with the BITBUCKET_USERNAME environment variable.
  # username = "me@example.com"

  # Bitbucket API token, created at https://id.atlassian.com/manage-profile/security/api-tokens.
  # Can also be set with the BITBUCKET_PASSWORD environment variable.
  # password = "ATATT3xFfGF0FakeToken"

  # Base URL of your Bitbucket Server.
  # Defaults to "https://api.bitbucket.org/2.0".
  # Can also be set with the BITBUCKET_API_BASE_URL environment variable.
  # base_url = "https://api.bitbucket.org/2.0"
}
```


