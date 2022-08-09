# Table: freshservice_asset

Obtain information about Assets from the FreshService instance.

## Examples

### List all assets

```sql
select
  *
from
  freshservice_asset;
```

### Obtain a specific asset, it's type and assignee

```sql
select
  a.display_id,
  a.name,
  a.assigned_on,
  t.name as asset_type,
  t.description,
  u.first_name,
  u.last_name,
  u.email
from
  freshservice.freshservice_asset a
inner join
  freshservice_asset_type t
on
  a.asset_type_id = t.id
inner join
  freshservice_requester u
on
  a.user_id = u.id
where
  a.id = 27001020436;
```
