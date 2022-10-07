package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableVendor() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_vendor",
		Description: "Obtain information on Vendors stored in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listVendors,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVendor,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: vendorColumns(),
	}
}

func vendorColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the vendor.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the vendor.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the vendor.",
			Type:        proto.ColumnType_STRING,
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
			Description: "Zip (postal) code of the location.",
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
			Description: "Timestamp at which the vendor was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the vendor was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getVendor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_vendor.getVendor", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	vendor, _, err := client.Vendors.GetVendor(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_vendor.getVendor", "query_error", err)
		return nil, fmt.Errorf("unable to obtain vendor with id %d: %v", id, err)
	}

	return vendor, nil
}

func listVendors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_vendor.listVendors", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListVendorsOptions{
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
		vendors, res, err := client.Vendors.ListVendors(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_vendor.listVendors", "query_error", err)
			return nil, fmt.Errorf("unable to obtain vendors: %v", err)
		}

		for _, vendor := range vendors.Collection {
			d.StreamListItem(ctx, vendor)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
