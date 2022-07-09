package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tablePurchaseOrder() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_purchase_order",
		Description: "Obtain information on Purchase Orders raised in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listPurchaseOrders,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getPurchaseOrder,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: purchaseOrderColumns(),
	}
}

func purchaseOrderColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the purchase order.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Title of the purchase order.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "po_number",
			Description: "Unique purchase order number.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "vendor_id",
			Description: "ID of the vendor associated with the purchase order.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "vendor_details",
			Description: "Details of the vendor with whom the order is placed.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "expected_delivery_date",
			Description: "Timestamp when delivery is expected for the purchase order.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "status",
			Description: "Status of the purchase order.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "shipping_address",
			Description: "Address to which the order should be shipped.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "billing_address",
			Description: "Address to which the order should be billed.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "billing_same_as_shipping",
			Description: "Set to true if billing address is same as shipping address.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_by",
			Description: "User ID of the agent who created the purchase order.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "department_id",
			Description: "ID of the department.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "currency_code",
			Description: "Currency unit used in the transaction.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "conversion_rate",
			Description: "Conversion rate to convert selected currency unit to helpdesk currency.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "discount_percentage",
			Description: "Percentage of discount on the order.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "tax_percentage",
			Description: "Percentage of tax on the order.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "shopping_cost",
			Description: "Total cost of shipping the order.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "purchase_items",
			Description: "Items to be ordered.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the purchase order was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the purchase order was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getPurchaseOrder(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_purchase_order.getPurchaseOrder", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	purchaseOrder, _, err := client.PurchaseOrders.GetPurchaseOrder(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_purchase_order.getPurchaseOrder", "query_error", err)
		return nil, fmt.Errorf("unable to obtain purchase order with id %d: %v", id, err)
	}

	return purchaseOrder, nil
}

func listPurchaseOrders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_purchase_order.listPurchaseOrders", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListPurchaseOrdersOptions{
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
		purchaseOrders, res, err := client.PurchaseOrders.ListPurchaseOrders(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_purchase_order.listPurchaseOrders", "query_error", err)
			return nil, fmt.Errorf("unable to obtain purchase orders: %v", err)
		}

		for _, purchaseOrder := range purchaseOrders.Collection {
			d.StreamListItem(ctx, purchaseOrder)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
