# Table: freshservice_software_user

Obtain information Users assigned to Software.

## Examples

### List all users of a specific software

```sql
select 
  *
from
  freshservice_software_user
where
  software_id = 20585369;
```

### Get a specific user for a software

```sql
select
  *
from
  freshservice_software_user
where
  software_id = 20585369
  and id = 1008564;
```
