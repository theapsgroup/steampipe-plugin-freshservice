package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableContract() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_contract",
		Description: "Obtain information about Contracts from the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listContracts,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getContract,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: contractColumns(),
	}
}

func contractColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the contract.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the contract.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the contract.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "vendor_id",
			Description: "ID of the vendor.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "auto_renew",
			Description: "True if the contract renews automatically.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "notify_expiry",
			Description: "True if the expiration notifications are configured for the contract.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "notify_before",
			Description: "Number of days before contract expiry date when the expiry notifications need to be sent.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "approver_id",
			Description: "ID of the agent who needs to approve the contract.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "start_date",
			Description: "Start date of the contract.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "end_date",
			Description: "End date of the contract.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "cost",
			Description: "Cost of the contract.",
			Type:        proto.ColumnType_STRING, // Formerly proto.ColumnType_DOUBLE but costs are to be strings (https://steampipe.io/docs/develop/plugin-release-checklist#column-types)
		},
		{
			Name:        "status",
			Description: "Status of the contract.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "contract_number",
			Description: "Unique reference number for the contract.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "contract_type_id",
			Description: "ID of the contract type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "visible_to_id",
			Description: "ID of agent group in FreshService to control visibility of the contract.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "notify_to",
			Description: "Array of email address whom should be notified of contract expiry.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "expiry_notified",
			Description: "True if the contract expiration notification has been sent.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "requester_id",
			Description: "ID of user whom created/renewed the contract.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "delegatee_id",
			Description: "ID of the agent whom the contract approval is delegated to",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the contract was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the contract was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Note: Can add the additional columns at later date for software specific contracts:
// software_id, license_type, billing_cycle, license_key, item_cost_details

// Hydrate Functions
func getContract(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.EqualsQuals["display_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_contract.getContract", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	contract, _, err := client.Contracts.GetContract(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_contract.getContract", "query_error", err)
		return nil, fmt.Errorf("unable to obtain contract with id %d: %v", id, err)
	}

	return contract, nil
}

func listContracts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_contract.listContracts", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListContractsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(30) {
			filter.PerPage = int(*limit)
		}
	}

	for {
		contracts, res, err := client.Contracts.ListContracts(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_contract.listContracts", "query_error", err)
			return nil, fmt.Errorf("unable to obtain contracts: %v", err)
		}

		for _, contract := range contracts.Collection {
			d.StreamListItem(ctx, contract)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
