# Table: freshservice_release_task

Allows for obtaining information on tasks associated to a specific Release.

You **MUST** specify a `release_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all tasks on a specific Release

```sql
select
  *
from
  freshservice_release_task
where
  release_id = 2011111111;
```
