package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableBusinessHour() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_business_hour",
		Description: "Retrieve Business Hours configuration of the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listBusinessHours,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getBusinessHours,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: businessHoursColumns(),
	}
}

func businessHoursColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the business hours configuration.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the business hours configuration.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description about the business hours configuration.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "is_default",
			Description: "True if the business hours configuration is the default present in Freshservice.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "time_zone",
			Description: "Time zone that the business hours configuration functions in",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "service_desk_hours",
			Description: "Contains the time at which the workday begins and ends for the seven days of the week.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "list_of_holidays",
			Description: "Contains the list of dates and names of holidays for the year.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the business hours configuration was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the business hours configuration was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getBusinessHours(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	bh, _, err := client.BusinessHours.GetBusinessHours(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain business hours configuration with id %d: %v", id, err)
	}

	return bh, nil
}

func listBusinessHours(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListBusinessHoursOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		bhs, res, err := client.BusinessHours.ListBusinessHours(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain business hours configurations: %v", err)
		}

		for _, bh := range bhs.Collection {
			d.StreamListItem(ctx, bh)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}