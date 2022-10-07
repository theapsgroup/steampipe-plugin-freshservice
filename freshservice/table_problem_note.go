package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableProblemNote() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_problem_note",
		Description: "Obtain notes for a specific Problem in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listProblemNotes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "problem_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: problemNoteColumns(),
	}
}

func problemNoteColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the problem note.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "user_id",
			Description: "User ID of the user who created the note.",
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
			Description: "Timestamp at which the note was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "problem_id",
			Description: "ID of the problem this note belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("problem_id"),
		},
	}
}

// Hydrate Functions
func listProblemNotes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	problemId := int(d.KeyColumnQuals["problem_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem_note.listProblemNotes", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	notes, _, err := client.Problems.ListProblemNotes(problemId)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem_note.listProblemNotes", "query_error", err)
		return nil, fmt.Errorf("unable to obtain problem notes: %v", err)
	}

	for _, note := range notes.Collection {
		d.StreamListItem(ctx, note)
	}

	return nil, nil
}
