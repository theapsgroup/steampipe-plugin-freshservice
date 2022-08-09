# Table: freshservice_release

Obtain information about Releases within the FreshService instance.

## Examples

### List all releases

```sql
select
  *
from
  freshservice_release;
```

### List high priority yet not completed releases

```sql
select
  r.id,
  r.subject,
  a.email as assignee,
  r.priority_desc,
  r.status_desc,
  r.planned_start_date,
  r.planned_end_date
from
  freshservice_release r
  left join freshservice_agent a on r.agent_id = a.id
where
  r.status < 5
  and r.priority >=3;
```
