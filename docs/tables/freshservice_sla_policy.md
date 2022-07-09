# Table: freshservice_sla_policy

Obtain information on Service Level Agreement Policies defined in the FreshService instance.

## Examples

### List all SLA Policies

```sql
select
  *
from
  freshservice_sla_policy;
```

### Get default SLA Policies

```sql
select
  id,
  name,
  active,
  deleted,
  category,
  sub_category,
  item_category
from
  freshservice_sla_policy
where
  is_default = true;
```
