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

### List all notes of all problems

```sql
select
  p.id,
  p.description_text,
  p.priority_desc,
  n.id as note_id,
  n.body_text as note
from
  freshservice_problem p
  left join freshservice_problem_note n on p.id = n.problem_id;
```
