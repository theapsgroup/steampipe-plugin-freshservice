package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableAssetType() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_asset_type",
		Description: "Obtain information about Asset Types from your FreshService instance.",
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
			Description: "ID of the asset type.",
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
			Description: "True if the asset type is visible.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the asset type was created",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the asset type was last updated",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getAssetType(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_asset_type.getAssetType", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	assetType, _, err := client.Assets.GetAssetType(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_asset_type.getAssetType", "query_error", err)
		return nil, fmt.Errorf("unable to obtain asset type with id %d: %v", id, err)
	}

	return assetType, nil
}

func listAssetTypes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_asset_type.listAssetTypes", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListAssetTypesOptions{
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
		assetTypes, res, err := client.Assets.ListAssetTypes(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_asset_type.listAssetTypes", "query_error", err)
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
