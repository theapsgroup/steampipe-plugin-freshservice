# Table: freshservice_problem

Obtain information about Problems in the FreshService instance.

## Examples

### List all problems

```sql
select
  *
from
  freshservice_problem;
```

### List problems that're not known issues

```sql
select
  *
from
  freshservice_problem
where
  known_error = false;
```
