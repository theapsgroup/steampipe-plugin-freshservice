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

### List software installations associated with a specific user (requester)

```sql
select
  *
from
  freshservice_software_installation fsi 
inner join
  freshservice_software fs 
on 
  fsi.software_id = fs.id
where
  fsi.user_id = 27000123;
```
