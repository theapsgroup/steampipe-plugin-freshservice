# Table: freshservice_change_note

Obtain information about Notes attached to Changes in the FreshService instance.

You **MUST** specify a `change_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all change notes for a specific change

```sql
select 
  *
from
  freshservice_change_note
where
  change_id = 12345;
```
