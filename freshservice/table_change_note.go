package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableChangeNote() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_change_note",
		Description: "Obtain notes for a specific Change in the FreshService Instance.",
		List: &plugin.ListConfig{
			Hydrate: listChangeNotes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "change_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: changeNoteColumns(),
	}
}

func changeNoteColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the change note.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "user_id",
			Description: "ID of the user who created the note.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "body",
			Description: "The body of the note in HTML format.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "body_text",
			Description: "The body of the note in plain text format.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "notify_emails",
			Description: "Array of addresses to which notifications are sent.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the note was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the note was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "change_id",
			Description: "Unique ID of the Change this note belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("change_id"),
		},
	}
}

// Hydrate Functions
func listChangeNotes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	changeId := int(d.KeyColumnQuals["change_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	notes, _, err := client.Changes.ListChangeNotes(changeId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain change notes: %v", err)
	}

	for _, note := range notes.Collection {
		d.StreamListItem(ctx, note)
	}

	return nil, nil
}
