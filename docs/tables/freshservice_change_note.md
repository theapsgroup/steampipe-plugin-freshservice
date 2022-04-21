# Table: freshservice_change_note

Obtain information about Notes attached to Changes in the FreshService instance.

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
