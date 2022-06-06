# Table: freshservice_contract

Obtain information about Contracts from the FreshService instance.

## Examples

### List all contracts

```sql
select
  *
from
 freshservice_contract;
```

### List contracts which did not yet start

```sql
select
  c.id,
  c.name,
  c.description,
  c.contract_number
  v.name as vendor,
  c.start_date
from
  freshservice_contract c
inner join
  freshservice_vendor v
on c.vendor_id = v.id
where
  c.start_date > NOW()::timestamp;
```
