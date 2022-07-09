# Table: freshservice_problem_timeentry

Obtain time entries for a specific Problem in the FreshService instance.

You **MUST** specify a `problem_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all time entries for a specific Problem

```sql
select
  *
from
  freshservice_problem_timeentry
where
  problem_id = 2011111111;
```
