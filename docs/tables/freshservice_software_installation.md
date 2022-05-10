# Table: freshservice_software_installation

Obtain information about Installations of Software registered in the FreshService instance.

## Examples

### List all installations for a specific software

```sql
select
  *
from
  freshservice_software_installation
where
  software_id = 465465131;
```
