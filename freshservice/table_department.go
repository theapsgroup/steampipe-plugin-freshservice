package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableDepartment() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_department",
		Description: "Information about Departments stored within the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listDepartments,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getDepartment,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: departmentColumns(),
	}
}

func departmentColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the department.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the department.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description about the department.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "head_user_id",
			Description: "Unique identifier of the agent or requester who serves as the head of the department.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "prime_user_id",
			Description: "Unique identifier of the agent or requester who serves as the prime user of the department.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "domains",
			Description: "Email domains associated with the department.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the department was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the department was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getDepartment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	department, _, err := client.Departments.GetDepartment(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain department with id %d: %v", id, err)
	}

	return department, nil
}

func listDepartments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListDepartmentsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		departments, res, err := client.Departments.ListDepartments(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain departments: %v", err)
		}

		for _, department := range departments.Collection {
			d.StreamListItem(ctx, department)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
