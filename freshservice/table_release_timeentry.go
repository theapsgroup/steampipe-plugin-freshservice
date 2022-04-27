package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableReleaseTimeEntry() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_release_timeentry",
		Description: "Obtain time entries for a specific Release",
		List: &plugin.ListConfig{
			Hydrate: listReleaseTimeEntries,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "release_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: releaseTimeEntryColumns(),
	}
}

func releaseTimeEntryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the time entry.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "start_time",
			Description: "The time at which the time entry is added. If a timer, which is in stopped state, is started again, this holds date_time at which the timer is started again.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "executed_at",
			Description: "Time at which the timer is executed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "timer_running",
			Description: "Set to true if timer is currently running.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "billable",
			Description: "Set as true if the time entry is billable.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "time_spent",
			Description: "The total amount of time spent by the timer in hh::mm format. This field cannot be set if timer_running is true.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "task_id",
			Description: "ID of the task assigned to the time entry.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "ID of the user/agent to whom this time entry is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "note",
			Description: "Description of the time entry.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "created_at",
			Description: "Time at which this time entry is created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Time at which the time entry is updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "release_id",
			Description: "ID of the Release the time entry belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("release_id"),
		},
	}
}

// Hydrate Functions
func listReleaseTimeEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	releaseId := int(d.KeyColumnQuals["release_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	entries, _, err := client.Releases.ListTimeEntries(releaseId)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain time entries: %v", err)
	}

	for _, entry := range entries.Collection {
		d.StreamListItem(ctx, entry)
	}

	return nil, nil
}
