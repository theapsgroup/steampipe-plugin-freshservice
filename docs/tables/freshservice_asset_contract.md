# Table: freshservice_asset_contract

Allows for obtaining information about contracts associated to an Asset.

## Examples

### List all contracts associated to a specific asset

```sql
select *
from freshservice_asset_contract
where asset_display_id = 1234;
```

### Obtain only active contracts for a specific asset

```sql
select *
from freshservice_asset_contract
where asset_display_id = 1234
and contract_status = 'Active';
```

> TODO: Insert example when Contracts table is added to link the two.
