# Table: freshservice_ticket_task

Allows for obtaining information on tasks associated to a specific Ticket.

You **MUST** specify a `ticket_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all tasks on a specific Ticket

```sql
select
  *
from
  freshservice_ticket_task
where
  ticket_id = 2010101010;
```
