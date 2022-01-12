# Table: freshservice_agent

Obtain information about Agents from your FreshService instance.

## Examples

### List all agents

```sql
select
  *
from
  freshservice_agent;
```

### List agents that have logged in within last 10 days

```sql
select
  *
from
  freshservice_agent
where
  last_login_at > now()::date-10;
```

### List agents and reporting manager name

```sql
select 
  a.id,
  concat(a.first_name, ' ', a.last_name) as agent_name,
  a.email,
  a.active,
  concat(m.first_name, ' ', m.last_name) as manager_name 
from 
  freshservice_agent a 
left join
  freshservice_agent m 
on 
  a.reporting_manager_id = m.id;
```
