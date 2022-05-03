package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableService() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_service",
		Description: "Obtains information about Service Items from the Service Catalog on the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listServiceItems,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getServiceItem,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: serviceColumns(),
	}
}

func serviceColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the service item.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the service item.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the service item.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "short_description",
			Description: "Short description of the service item.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "delivery_time",
			Description: "Estimated delivery time of the service item (in hours).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "display_id",
			Description: "ID of the service item specific to your account.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "category_id",
			Description: "ID of the category of the service item.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "product_id",
			Description: "The ID of the product mapped to the item.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "quantity",
			Description: "The quantity set against the service item.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "item_type",
			Description: "1 indicates a normal item. 2 indicates a loaner item.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "ci_type_id",
			Description: "ID of the asset type associated with the product.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "visibility",
			Description: "1 denotes draft and 2 denotes published.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "group_visibility",
			Description: "1 denotes visibility to all requesters. 2 for restricted visibility.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "deleted",
			Description: "Set to true if the service item is deleted.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "cost_visibility",
			Description: "Set to true if cost should be visible to the requester.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "delivery_time_visibility",
			Description: "Set to true if delivery time of the service item should be visible to the requester.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "botified",
			Description: "Set to true if the service item is 'bot-ready'.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "allow_attachments",
			Description: "Set to true if the requester is allowed to attach a file.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "allow_quantity",
			Description: "Set to true if the requester is allowed request more than 1 quantity.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "is_bundle",
			Description: "Set to true if the service item contains child items.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "create_child",
			Description: "Set to true if child items can be created as separate service requests.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "cost",
			Description: "Cost of the service item.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "created_at",
			Description: "The time at which the service item was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "The time at which the service item was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

/*
type ServiceItem struct {
	Cost                   float32   `json:"cost"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
*/

// Hydrate Functions
func getServiceItem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	service, _, err := client.Services.GetServiceItem(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain service item with id %d: %v", id, err)
	}

	return service, nil
}

func listServiceItems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	// Note: List operation returns all results without pagination.
	serviceItems, _, err := client.Services.ListServiceItems()
	if err != nil {
		return nil, fmt.Errorf("unable to obtain service items: %v", err)
	}

	for _, serviceItem := range serviceItems.Collection {
		d.StreamListItem(ctx, serviceItem)
	}

	return nil, nil
}
