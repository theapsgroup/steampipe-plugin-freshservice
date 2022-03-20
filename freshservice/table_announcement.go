package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableAnnouncement() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_announcement",
		Description: "FreshService Announcements",
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
			Description: "Unique identifier of the Announcement",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_by",
			Description: "Unique identifier of the agent that created this Announcement",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "state",
			Description: "State of the Announcement - active, archived, scheduled",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "title",
			Description: "Title of the Announcement",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "body",
			Description: "Body of the Announcement",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "body_html",
			Description: "Body of the Announcement in HTML format",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "visible_from",
			Description: "Timestamp at which Announcement becomes active",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "visible_to",
			Description: "Timestamp until which Announcement is active",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "visibility",
			Description: "Who can see the announcement? Values - everyone, agents_only, agents_and_groups",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "departments",
			Description: "Array of Department IDs that can view this Announcement",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "groups",
			Description: "Array of Group IDs that can view this Announcement",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "is_read",
			Description: "True if the logged-in-user has read the announcement",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "send_email",
			Description: "True if the announcement needs to be sent via email as well. False, otherwise",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "additional_emails",
			Description: "Additional email addresses to which the announcement needs to be sent",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Announcement creation timestamp",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Announcement updated timestamp",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getAnnouncement(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	announcement, _, err := client.Announcements.GetAnnouncement(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain announcement with id %d: %v", id, err)
	}

	return announcement, nil
}

func listAnnouncements(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListAnnouncementsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	q := d.KeyColumnQuals
	if q["state"] != nil {
		filter.State = q["state"].GetStringValue()
	}

	for {
		announcements, res, err := client.Announcements.ListAnnouncements(filter)
		if err != nil {
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
