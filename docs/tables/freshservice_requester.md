# Table: freshservice_requester

Obtain information about Requesters (users) from your FreshService instance.

## Examples

### List all requesters

```sql
select
  *
from
  freshservice_requester;
```

### List users and reporting manager

```sql
select
  r.id,
  concat(r.first_name, ' ', r.last_name) as requester_name,
  r.email,
  r.active,
  concat(m.first_name, ' ', m.last_name) as manager_name
from
  freshservice_requester r
  left join freshservice_requester m on r.reporting_manager_id = m.id;
```

### List inactive users

```sql
select
  *
from
  freshservice_requester
where
  active = false;
```
