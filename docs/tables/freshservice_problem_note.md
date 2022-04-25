# Table: freshservice_problem_note

Obtain information about Notes attached to Problems in the FreshService instance.

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
