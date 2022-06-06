package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableBusinessHour() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_business_hour",
		Description: "Obtain the Business Hours configurations from the FreshService instance.",
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
			Description: "ID of the business hours configuration.",
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
			Description: "Time zone that the business hours configuration functions in.",
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
			Description: "Timestamp when the business hours configuration were created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the business hours configuration were last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getBusinessHours(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_business_hour.getBusinessHours", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	bh, _, err := client.BusinessHours.GetBusinessHours(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_business_hour.getBusinessHours", "query_error", err)
		return nil, fmt.Errorf("unable to obtain business hours configuration with id %d: %v", id, err)
	}

	return bh, nil
}

func listBusinessHours(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_business_hour.listBusinessHours", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListBusinessHoursOptions{
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
		bhs, res, err := client.BusinessHours.ListBusinessHours(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_business_hour.listBusinessHours", "query_error", err)
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
