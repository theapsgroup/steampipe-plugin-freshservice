# Table: freshservice_service

Obtain information about Service Items from the Service Catalog on the FreshService instance.

## Examples

### List all service items

```sql
select
  *
from
  freshservice_service;
```

### Get a specific service item by id

```sql
select
  name,
  short_description
from
  freshservice_service
where
  id = 27000112233;
```
