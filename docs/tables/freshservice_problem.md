# Table: freshservice_problem

Obtain information on Problems raised in the FreshService instance.

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

### List problems assigned to a specific agent

```sql
select
  id,
  description
from
  freshservice_problem
where
  agent_id = 2578963125;
```
