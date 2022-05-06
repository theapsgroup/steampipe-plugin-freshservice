package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableAssetType() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_asset_type",
		Description: "Information about Asset Types in FreshService",
		List: &plugin.ListConfig{
			Hydrate: listAssetTypes,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAssetType,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: assetTypeColumns(),
	}
}

func assetTypeColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the asset type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the asset type.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Short description of the asset type.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "parent_asset_type_id",
			Description: "ID of the parent asset type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "visible",
			Description: "Visibility of the default asset type. Set to true if the asset type is visible.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the asset type was created",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the asset type was last updated",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getAssetType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	assetType, _, err := client.Assets.GetAssetType(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain asset type with id %d: %v", id, err)
	}

	return assetType, nil
}

func listAssetTypes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListAssetTypesOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		assetTypes, res, err := client.Assets.ListAssetTypes(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain asset types: %v", err)
		}

		for _, assetType := range assetTypes.Collection {
			d.StreamListItem(ctx, assetType)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}
	return nil, nil
}
