# Table: freshservice_ticket_timeentry

Allows for obtaining all the time entries against a specific Ticket.

You **MUST** specify a `ticket_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all time entries for a given Ticket

```sql
select
  *
from
  freshservice_ticket_timeentry
where
  ticket_id = 2010101010;
```

### List all non-billable time entries for a given Ticket

```sql
select
  *
from
  freshservice_ticket_timeentry
where
  ticket_id = 2010101010
and 
  billable = false;
```
