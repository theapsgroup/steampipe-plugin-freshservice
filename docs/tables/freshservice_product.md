# Table: freshservice_product

Obtain information on Products in the FreshService instance.

## Examples

### Get a specific Product by id

```sql
select 
  *
from
  freshservice_product
where
  id = 27000001;
```

### List all Products

```sql
select
  *
from
  freshservice_product;
```

### List all Products for visible Asset Types

```sql
select
  p.id,
  p.name,
  p.status,
  p.manufacturer
from
  freshservice_product p
  inner join freshservice_asset_type at on p.asset_type_id = at.id
where
  at.visible = true;
```
