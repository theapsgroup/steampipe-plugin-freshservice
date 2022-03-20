# Table: freshservice_asset_type

Obtain information about Asset Types from your FreshService instance.

## Examples

### List all asset types

```sql
select
  *
from
  freshservice_asset_type;
```

### List all child asset types of a parent asset type

```sql
select
  *
from
  freshservice_asset_type
where
    parent_asset_type_id = 20069004;
```

### List all assets of a specific asset type

```sql
select
  a.*
from
  freshservice_asset a 
left join
  freshservice_aset_type t 
on a.asset_type_id = t.id
and t.name = 'MY-TYPE';
```
