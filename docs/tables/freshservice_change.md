# Table: freshservice_change

Obtain information about Changes in the FreshService instance.

### List all changes

```sql
select
  *
from
  freshservice_change;
```

### Obtain changes assigned to a specific agent

```sql
select
  c.id as change_id,
  c.description_text as change,
  c.status_desc as status,
  a.email as agent
from
  freshservice_change c
  inner join freshservice_agent a on c.agent_id = a.id
where
  a.id = 12345;
```
