# Table: freshservice_contract_type

Obtain information about Contract Types in the FreshService instance.

## Examples

### List all contract types

```sql
select
  *
from
  freshservice_contract_type;
```

### List all contracts of a specific type

```sql
select
  c.id,
  c.name,
  c.description,
  c.contract_number,
  t.name as contract_type,
  c.start_date
from
  freshservice_contract c
  inner join freshservice_contract_type t on c.contract_type_id = t.id
where
  t.id = 20009451249;
```
