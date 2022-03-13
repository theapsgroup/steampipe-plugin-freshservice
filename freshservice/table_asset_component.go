package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAssetComponent() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_asset_component",
		Description: "Information about the Components of a specific Asset.",
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
			Description: "Unique ID of the component.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "component_type",
			Description: "Type of the Component. (Example: Processor, Memory)",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "component_data",
			Description: "Details of the Component.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the component was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the component was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listAssetComponents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	q := d.KeyColumnQuals
	displayId := int(q["asset_display_id"].GetInt64Value())

	if displayId == 0 {
		return nil, fmt.Errorf("freshservice_asset_component List call requires an '=' qualifier for 'asset_display_id'")
	}

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	components, _, err := client.Assets.ListAssetComponents(displayId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain asset components: %v", err)
	}

	for _, component := range components.Collection {
		d.StreamListItem(ctx, component)
	}

	return nil, nil
}
