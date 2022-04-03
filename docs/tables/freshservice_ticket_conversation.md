# Table: freshservice_ticket_conversation

Allows for obtaining information on conversations for a specific Ticket.

You **MUST** specify a `ticket_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all conversations on a specific Ticket

```sql
select
  *
from
  freshservice_ticket_conversation
where
  ticket_id = 2010101010;
```
