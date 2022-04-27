# Table: freshservice_problem_timeentry

Allows for obtaining information on time entries associated to a specific Problem.

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
