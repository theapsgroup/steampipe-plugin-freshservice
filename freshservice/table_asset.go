package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableAsset() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_asset",
		Description: "Information about Assets stored within the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listAssets,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAsset,
			KeyColumns: plugin.SingleColumn("display_id"),
		},
		Columns: assetColumns(),
	}
}

func assetColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the asset.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "display_id",
			Description: "Display ID of the asset.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the asset.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the asset.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "asset_type_id",
			Description: "ID of the asset type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "asset_tag",
			Description: "Asset tag of the asset.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "impact",
			Description: "Impact of the asset.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "author_type",
			Description: "Indicates whether the asset was created by a user or discovery tools (Probe or Agent).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "usage_type",
			Description: "Usage type of the asset (Loaner / Permanent).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "user_id",
			Description: "ID of the associated user (Used By).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "location_id",
			Description: "ID of the associated location.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "department_id",
			Description: "ID of the associated department.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "ID of the associated agent (Managed By).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "group_id",
			Description: "ID of the associated agent group (Managed By Group).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "assigned_on",
			Description: "Date and time when the asset was assigned.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the asset was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the asset was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getAsset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["display_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	asset, _, err := client.Assets.GetAsset(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain asset with display_id %d: %v", id, err)
	}

	return asset, nil
}

func listAssets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListAssetsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		agents, res, err := client.Assets.ListAssets(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain assets: %v", err)
		}

		for _, agent := range agents.Collection {
			d.StreamListItem(ctx, agent)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
