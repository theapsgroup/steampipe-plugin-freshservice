# Table: freshservice_problem_note

Obtain information about Notes attached to Problems in the FreshService instance.

You **MUST** specify a `problem_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all notes for a specific problem

```sql
select 
  *
from
  freshservice_problem_note
where
  problem_id = 12345;
```
