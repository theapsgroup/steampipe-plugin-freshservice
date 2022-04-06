package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableProduct() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_product",
		Description: "Obtain information on Products in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listProducts,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getProduct,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: productColumns(),
	}
}

func productColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the product.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the Product.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "HTML content of the product.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description_text",
			Description: "Description of the product in plain text.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "asset_type_id",
			Description: "ID of the asset type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "manufacturer",
			Description: "Manufacturer of the product",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "status",
			Description: "Status of the product (In Production, In Pipeline, Retired).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "mode_of_procurement",
			Description: "Mode of procurement of the product (Buy, Lease, Both).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "depreciation_type_id",
			Description: "ID of the depreciation type.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the product was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the product was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getProduct(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	product, _, err := client.Products.GetProduct(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain product with id %d: %v", id, err)
	}

	return product, nil
}

func listProducts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListProductsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		products, res, err := client.Products.ListProducts(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain products: %v", err)
		}

		for _, product := range products.Collection {
			d.StreamListItem(ctx, product)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
