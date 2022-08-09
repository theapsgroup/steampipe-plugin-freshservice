# Table: freshservice_solution_article

Obtain information on Solution Articles stored in the FreshService instance.

## Examples

### List all solution articles

```sql
select
  *
from
  freshservice_solution_article;
```

### Get a specific solution article by id

```sql
select
  id,
  title,
  status
from
  freshservice_solution_article
where
  id = 269696969;
```

### List all solution articles for a specific solution category

```sql
select
  a.id,
  a.title,
  a.status,
  c.name as category,
  a.views,
  a.review_date
from
  freshservice_solution_article a
  inner join freshservice_solution_category c on a.category_id = c.id;
```
