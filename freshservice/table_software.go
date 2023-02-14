package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableSoftware() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_software",
		Description: "Obtain information on Software stored in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listSoftware,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getSoftware,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: softwareColumns(),
	}
}

func softwareColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the software.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "application_type",
			Description: "Type of the software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "status",
			Description: "Status of the software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "category",
			Description: "Category of the software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "notes",
			Description: "Notes about the software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "user_count",
			Description: "Number of users using this software.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "installation_count",
			Description: "Number of devices the software is installed on.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "publisher_id",
			Description: "ID of the publisher.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "managed_by_id",
			Description: "User ID of the agent managing the software.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when software record is created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the software record was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getSoftware(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.EqualsQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_software.getSoftware", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	software, _, err := client.Software.GetApplication(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_software.getSoftware", "query_error", err)
		return nil, fmt.Errorf("unable to obtain software with id %d: %v", id, err)
	}

	return software, nil
}

func listSoftware(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_software.listSoftware", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListApplicationsOptions{
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
		allSoftware, res, err := client.Software.ListApplications(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_software.listSoftware", "query_error", err)
			return nil, fmt.Errorf("unable to obtain software: %v", err)
		}

		for _, software := range allSoftware.Collection {
			d.StreamListItem(ctx, software)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
