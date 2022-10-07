package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableReleaseNote() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_release_note",
		Description: "Obtain notes for a specific Release in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listReleaseNotes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "release_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: releaseNoteColumns(),
	}
}

func releaseNoteColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the release note.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "user_id",
			Description: "User ID of who created the note.",
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
			Name:        "release_id",
			Description: "ID of the release this note belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("release_id"),
		},
	}
}

// Hydrate Functions
func listReleaseNotes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	releaseId := int(d.KeyColumnQuals["release_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_release_note.listReleaseNotes", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	notes, _, err := client.Releases.ListReleaseNotes(releaseId)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_release_note.listReleaseNotes", "query_error", err)
		return nil, fmt.Errorf("unable to obtain release notes: %v", err)
	}

	for _, note := range notes.Collection {
		d.StreamListItem(ctx, note)
	}

	return nil, nil
}
