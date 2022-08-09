# Table: freshservice_release_note

Obtain information about Notes attached to Release in the FreshService instance.

You **MUST** specify a `release_id` in the `WHERE` or `JOIN` clause.

## Examples

### List all notes for a specific release

```sql
select
  *
from
  freshservice_release_note
where
  release_id = 200012345;
```
