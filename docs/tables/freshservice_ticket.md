# Table: freshservice_ticket

Obtain information about Tickets in the FreshService instance.

## Examples

### List all tickets

```sql
select
  *
from
  freshservice_ticket;
```

### List tickets raised by a specific email address

```sql
select
  *
from
  freshservice_ticket
where
  email = 'some@email.here';
```

### List tickets that are `Open` and `Urgent`

```sql
select
  t.id,
  t.subject,
  t.name as requester,
  t.status_desc as status,
  t.priority_desc as priority,
  t.category,
  t.sub_category,
  t.item_category,
  t.due_by,
  concat(a.first_name, ' ', a.last_name) as agent
from
  freshservice_ticket t
  left outer join freshservice_agent a on t.responder_id = a.id
where
  status = 2
  and priority = 4;
```
