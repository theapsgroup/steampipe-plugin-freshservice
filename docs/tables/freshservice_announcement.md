# Table: freshservice_announcement

Obtain information about Announcements from the FreshService instance.

## Examples

### List all announcements

```sql
select
  *
from
  freshservice_announcement;
```

### List announcements scheduled for the future

```sql
select
  *
from
  freshservice_announcement
where
  state = 'scheduled';
```

### List announcements created by a specific agent

```sql
select
  a.title,
  a.body,
  concat(ag.first_name, ' ', ag.last_name) as agent
from
  freshservice_announcement a
  left join freshservice_agent ag on a.created_by = ag.id
where
  ag.email = 'example@agent.com';
```
