package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableAssetContract() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_asset_contract",
		Description: "Information on contracts for a specific Asset",
		List: &plugin.ListConfig{
			Hydrate:    listAssetContracts,
			KeyColumns: plugin.SingleColumn("asset_display_id"),
		},
		Columns: assetContractColumns(),
	}
}

func assetContractColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "asset_display_id",
			Description: "Display ID of the parent asset",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("asset_display_id"),
		},
		{
			Name:        "id",
			Description: "Contract ID specific to your account.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "contract_id",
			Description: "Unique Contract Number",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "contract_type",
			Description: "Type of the Contract. (Example: Lease, Maintenance)",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "contract_name",
			Description: "Subject/Title of the Contract",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "contract_status",
			Description: "Status of the contract.(Example: Active, Draft)",
			Type:        proto.ColumnType_STRING,
		},
	}
}

// Hydrate Functions
func listAssetContracts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.KeyColumnQuals
	displayId := int(q["asset_display_id"].GetInt64Value())

	if displayId == 0 {
		return nil, fmt.Errorf("freshservice_asset_contract List call requires an '=' qualifier for 'asset_display_id'")
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	contracts, _, err := client.Assets.ListAssetContracts(displayId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain asset contracts: %v", err)
	}

	for _, contract := range contracts.Collection {
		d.StreamListItem(ctx, contract)
	}

	return nil, nil
}
