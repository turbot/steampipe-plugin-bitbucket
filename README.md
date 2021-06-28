![image](https://hub.steampipe.io/images/plugins/turbot/bitbucket-social-graphic.png)

# Bitbucket Plugin for Steampipe

Use SQL to query repositories, projects, merge requests and more from Bitbucket.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/bitbucket)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/bitbucket/tables)
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-bitbucket/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install bitbucket
```

Run a query:

```sql
select name, uuid, is_private, full_name from bitbucket_my_repository;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone git@github.com:turbot/steampipe-plugin-bitbucket
cd steampipe-plugin-bitbucket
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/bitbucket.spc
```

Try it!

```
steampipe query
> .inspect bitbucket
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-bitbucket/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Bitbucket Plugin](https://github.com/turbot/steampipe-plugin-bitbucket/labels/help%20wanted)
