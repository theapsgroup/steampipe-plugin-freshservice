package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableContractType() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_contract_type",
		Description: "Obtain information about Contract Types in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listContractTypes,
		},
		Columns: contractTypeColumns(),
	}
}

func contractTypeColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the contract type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the contract type.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the contract type.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "needs_approval",
			Description: "True if the contract type needs approval.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "is_default",
			Description: "True if the contract type is a default (or custom) type.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the contract type was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the contract type was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listContractTypes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_contract_type.listContractTypes", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	contractTypes, _, err := client.Contracts.ListContractTypes()
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_contract_type.listContractTypes", "query_error", err)
		return nil, fmt.Errorf("unable to obtain contract types: %v", err)
	}

	for _, contractType := range contractTypes.Collection {
		d.StreamListItem(ctx, contractType)
	}

	return nil, nil
}
