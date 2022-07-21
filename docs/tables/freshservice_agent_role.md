# Table: freshservice_agent_role

Obtain information about Agent Roles from the FreshService instance.

## Examples

### List all agent roles

```sql
select
  *
from
  freshservice_agent_role;
```

### List all agents with their roles

```sql
select
   a.email,
   r.name as role
from 
  freshservice.freshservice_agent a
cross join 
  jsonb_to_recordset(a.roles) as ar(role_id bigint)
inner join
  freshservice.freshservice_agent_role r 
on 
  r.id::bigint = ar.role_id
order by 
  email;
```
