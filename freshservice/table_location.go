package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableLocation() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_location",
		Description: "Obtain information on locations registered in the FreshService Instance.",
		List: &plugin.ListConfig{
			Hydrate: listLocations,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getLocation,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: locationColumns(),
	}
}

func locationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the location.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the location.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "parent_location_id",
			Description: "ID of the parent location.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "primary_contact_id",
			Description: "User ID of the primary contact.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "line1",
			Description: "Address line 1.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Address.Line1"),
		},
		{
			Name:        "line2",
			Description: "Address line 2.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Address.Line2"),
		},
		{
			Name:        "city",
			Description: "Name of the city.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Address.City"),
		},
		{
			Name:        "state",
			Description: "Name of the state.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Address.State"),
		},
		{
			Name:        "zipcode",
			Description: "Zip/Postal code of the location.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Address.ZipCode"),
		},
		{
			Name:        "country",
			Description: "Name of the country.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Address.Country"),
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the location was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the location was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getLocation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_location.getLocation", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	location, _, err := client.Locations.GetLocation(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_location.getLocation", "query_error", err)
		return nil, fmt.Errorf("unable to obtain location with id %d: %v", id, err)
	}

	return location, nil
}

func listLocations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_location.listLocations", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListLocationsOptions{
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
		locations, res, err := client.Locations.ListLocations(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_location.listLocations", "query_error", err)
			return nil, fmt.Errorf("unable to obtain asset types: %v", err)
		}

		for _, location := range locations.Collection {
			d.StreamListItem(ctx, location)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}
	return nil, nil
}
