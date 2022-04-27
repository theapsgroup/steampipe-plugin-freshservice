# Table: freshservice_release_timeentry

Allows for obtaining information on time entries associated to a specific Release.

You **MUST** specify a `release_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all time entries for a specific Release

```sql
select
  *
from
  freshservice_release_timeentry
where
  release_id = 2011111111;
```
