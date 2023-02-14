package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAssetComponent() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_asset_component",
		Description: "Obtain information about the Components of a specific Asset.",
		List: &plugin.ListConfig{
			Hydrate:    listAssetComponents,
			KeyColumns: plugin.SingleColumn("asset_display_id"),
		},
		Columns: assetComponentColumns(),
	}
}

func assetComponentColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "asset_display_id",
			Description: "Display ID of the parent asset",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("asset_display_id"),
		},
		{
			Name:        "id",
			Description: "ID of the component.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "component_type",
			Description: "Type of the component. (Example: Processor, Memory)",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "component_data",
			Description: "Details of the component.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the component was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the component was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listAssetComponents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.EqualsQuals
	displayId := int(q["asset_display_id"].GetInt64Value())

	if displayId == 0 {
		err := fmt.Errorf("freshservice_asset_component List call requires an '=' qualifier for 'asset_display_id'")
		plugin.Logger(ctx).Error("freshservice_asset_component.listAssetComponents", "missing_qualifier_error", err)
		return nil, err
	}

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_asset_component.listAssetComponents", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	components, _, err := client.Assets.ListAssetComponents(displayId)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_asset_component.listAssetComponents", "query_error", err)
		return nil, fmt.Errorf("unable to obtain asset components: %v", err)
	}

	for _, component := range components.Collection {
		d.StreamListItem(ctx, component)
	}

	return nil, nil
}
