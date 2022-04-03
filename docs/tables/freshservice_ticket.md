# Table: freshservice_ticket

Obtain information about Tickets in the FreshService instance.

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
  *
from
  freshservice_ticket
where
  status = 2
and 
  priority = 4;
```
