package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableSolutionCategory() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_solution_category",
		Description: "Obtain information about Solution Categories in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listSolutionCategories,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: solutionCategoryColumns(),
	}
}

func solutionCategoryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the solution category.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the solution category.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description about the solution category.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "position",
			Description: "The rank of the solution category in the category listing.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "default_category",
			Description: "Set to true if the solution category is the default one.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "visible_in_portals",
			Description: "Array of portal IDs where this category is visible.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the solution category was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the solution category was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listSolutionCategories(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_solution_category.listSolutionCategories", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListSolutionCategoriesOptions{
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

	q := d.KeyColumnQuals

	if q["id"] != nil {
		catId := int(q["id"].GetInt64Value())

		category, _, err := client.Solutions.GetSolutionCategory(catId)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_solution_category.listSolutionCategories", "query_error", err)
			return nil, fmt.Errorf("unable to obtain solution category with id %d: %v", catId, err)
		}

		d.StreamListItem(ctx, category)
	} else {
		for {
			categories, res, err := client.Solutions.ListSolutionCategories(&filter)
			if err != nil {
				plugin.Logger(ctx).Error("freshservice_solution_category.listSolutionCategories", "query_error", err)
				return nil, fmt.Errorf("unable to obtain solution categories: %v", err)
			}

			for _, category := range categories.Collection {
				d.StreamListItem(ctx, category)
			}

			if res.Header.Get("link") == "" {
				break
			}

			filter.Page += 1
		}
	}

	return nil, nil
}
