---
organization: The APS Group
category: ["SaaS"]
icon_url: "/images/plugins/theapsgroup/freshservice.svg"
brand_color: "#148FCC"
display_name: "FreshService"
short_name: "freshservice"
description: "Steampipe plugin for querying FreshService agents, assets, tickets and other resources."
og_description: Query FreshService with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/theapsgroup/freshservice-social-graphic.png"
---

# FreshService + Steampipe

[FreshService](https://freshservice.com/) is a ITSM (IT Service Management) solution provided as a SaaS.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  id,
  name,
  active,
  category
from
  freshservice_sla_policy;
```

```
+-------------+------------------------------+--------+----------+
| id          | name                         | active | category |
+-------------+------------------------------+--------+----------+
| 97000001922 | Some Procurement SLA Policy  | true   | custom   |
| 97000012364 | IT Support SLA               | true   | custom   |
| 97000000124 | Default SLA Policy           | true   | default  |
+-------------+------------------------------+--------+----------+
```

## Documentation

- [Table definitions / examples](https://hub.steampipe.io/plugins/theapsgroup/freshservice/tables)

## Getting Started

### Installation

```shell
steampipe plugin install theapsgroup/freshservice
```

### Prerequisites

- FreshService account (this will give you a domain for your instance `https://domain.freshservice.com/`)
- API Token for an agent with admin permissions.

### Configuration

> Note: Configuration file will take precedence over Env Vars.

Configuration can be done via Environment Variables or via the Configuration file `~./steampipe/config/freshservice.spc`.

Environment Variables:
- `FRESHSERVICE_DOMAIN` : The friendly sub-domain at which your instance is deployed (example: `my-corp` if your instance is `https://my-corp.freshservice.com`).
- `FRESHSERVICE_TOKEN` : The API token you wish to use.

Configuration File:

```hcl
connection "freshservice" {
  plugin   = "theapsgroup/freshservice"
  domain   = "my-corp"
  token    = "34vt5394t534rv4tvr435v74b395t34qv9q"
}
```

### Testing

A quick test can be performed from your terminal with:

```shell
steampipe query "select * from freshservice_ticket"
```

## Get involved

* Open source: https://github.com/theapsgroup/steampipe-plugin-freshservice
* Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
