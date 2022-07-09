# Table: freshservice_purchase_order

Obtain information on Purchase Orders raised in the FreshService instance.

## Examples

### List all purchase orders

```sql
select
  *
from
  freshservice_purchase_order;
```

### List purchase orders from a specific vendor

```sql
select
  *
from
  freshservice_purchase_order
where
  vendor_id = 20591913;
```
