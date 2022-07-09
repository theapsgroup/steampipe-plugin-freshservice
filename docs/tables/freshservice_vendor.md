# Table: freshservice_vendor

Obtain information about Vendors stored in the FreshService instance.

## Examples

### List all vendors

```sql
select
  id,
  name,
  description
from
  freshservice_vendor;
```

### Get a specific vendor

```sql
select
  name,
  description,
  line1 as street,
  city,
  state,
  zipcode,
  country
from
  freshservice_vendor
where
  id = 296358410;
```
