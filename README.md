# FreshService plugin for Steampipe

Use SQL to query information including Tickets, Agents, Assets and more from FreshService.

- **[Get started →](https://hub.steampipe.io/plugins/theapsgroup/freshservice)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/theapsgroup/steampipe-plugin-freshservice/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install theapsgroup/freshservice
```

Setup the configuration:

```shell
vi ~/.steampipe/config/freshservice.spc
```

or set the following Environment Variables

- `FRESHSERVICE_DOMAIN : The friendly sub-domain at which your instance is deployed (example: `my-corp` if your instance is `https://my-corp.freshservice.com`)
- `FRESHSERVICE_TOKEN` : The API Key / Token to use.

Run a query:

```sql
select
  id,
  name,
  active,
  category
from
  freshservice_sla_policy;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)
- [FreshService](https://freshservice.com/)

Clone:

```sh
git clone https://github.com/theapsgroup/steampipe-plugin-freshservice.git
cd steampipe-plugin-freshservice
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/freshservice.spc
```

Try it!

```
steampipe query
> .inspect freshservice
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)
