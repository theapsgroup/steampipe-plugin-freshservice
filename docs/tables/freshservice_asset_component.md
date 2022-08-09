# Table: freshservice_asset_component

Allows for obtaining information about component parts of Assets from within the FreshService instance.

You **MUST** specify an `asset_display_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all components for a specific asset

```sql
select
  *
from
  freshervice_asset_component
where
  asset_display_id = 1234;
```

### List all components of all Assets

```sql
select
  a.id,
  a.display_id,
  a.name,
  a.asset_tag,
  c.id as component_id,
  c.component_type,
  c.component_data
from
  freshservice_asset a
  left join freshservice_asset_component c on a.display_id = c.asset_display_id;
```
