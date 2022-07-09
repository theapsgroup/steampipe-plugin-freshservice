package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableAnnouncement() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_announcement",
		Description: "Obtain information about Announcements from the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listAnnouncements,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "state",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAnnouncement,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: announcementColumns(),
	}
}

func announcementColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the announcement.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_by",
			Description: "ID of the agent that created this announcement.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "state",
			Description: "State of the announcement - (active, archived, scheduled).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "title",
			Description: "Title of the announcement.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "body",
			Description: "Body of the announcement.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "body_html",
			Description: "Body of the announcement in HTML format.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "visible_from",
			Description: "Timestamp at which the announcement becomes active.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "visible_to",
			Description: "Timestamp until the announcement is active.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "visibility",
			Description: "Visibility of the announcement - (everyone, agents_only, agents_and_groups).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "departments",
			Description: "Array of Department IDs that can view this announcement.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "groups",
			Description: "Array of Group IDs that can view this announcement.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "is_read",
			Description: "True if the logged-in-user has read the announcement.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "send_email",
			Description: "True if the announcement needs to be sent via email.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "additional_emails",
			Description: "Array of additional email addresses to which the announcement needs to be sent.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the announcement was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the announcement was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getAnnouncement(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_announcement.getAnnouncement", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	announcement, _, err := client.Announcements.GetAnnouncement(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_announcement.getAnnouncement", "query_error", err)
		return nil, fmt.Errorf("unable to obtain announcement with id %d: %v", id, err)
	}

	return announcement, nil
}

func listAnnouncements(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_announcement.listAnnouncements", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListAnnouncementsOptions{
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
	if q["state"] != nil {
		filter.State = q["state"].GetStringValue()
	}

	for {
		announcements, res, err := client.Announcements.ListAnnouncements(filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_announcement.listAnnouncements", "query_error", err)
			return nil, fmt.Errorf("unable to obtain announcements: %v", err)
		}

		for _, announcement := range announcements.Collection {
			d.StreamListItem(ctx, announcement)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}
	return nil, nil
}
