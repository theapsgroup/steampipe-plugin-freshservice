# Table: freshservice_problem_task

Allows for obtaining information on tasks associated to a specific Problem.

You **MUST** specify a `problem_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all tasks on a specific Problem

```sql
select
  *
from
  freshservice_problem_task
where
  problem_id = 2011111111;
```
