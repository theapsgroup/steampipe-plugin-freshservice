# Table: freshservice_software

Obtain information on Software stored in the FreshService instance.

## Examples

### List all software registered in FreshService

```sql
select
  *
from
  freshservice_software;
```

### List software managed by a specific agent

```sql
select
  s.id,
  s.name,
  s.category,
  s.user_count,
  concat(a.first_name, ' ', a.last_name) as agent_name
from
  freshservice_software s
  inner join freshservice_agent a on s.managed_by_id = a.id
where
  a.email = 'test@example.com';
```
