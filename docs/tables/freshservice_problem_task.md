# Table: freshservice_problem_task

Obtain tasks based on an associated Problem in the FreshService instance.

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

### List all overdue tasks of a specific problem

```sql
select
  *
from
  freshservice_problem_task
where
  problem_id = 2011111111
and
  due_date < NOW()::timestamp;
```
