package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableProblemTimeEntry() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_problem_timeentry",
		Description: "Obtain time entries for a specific Problem in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listProblemTimeEntries,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "problem_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: problemTimeEntryColumns(),
	}
}

func problemTimeEntryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the time entry.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "start_time",
			Description: "Timestamp when the time entry is added. If a timer, which is in stopped state, is started again, this holds date_time at which the timer is started again.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "executed_at",
			Description: "Timestamp when the timer is executed.",
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
			Description: "The total amount of time spent by the timer in hh::mm format.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "task_id",
			Description: "ID of the task associated with the time entry.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "User ID of the agent to whom this time entry is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "note",
			Description: "Description of the time entry.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the time entry is created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the time entry was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "problem_id",
			Description: "ID of the problem the time entry belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("problem_id"),
		},
	}
}

// Hydrate Functions
func listProblemTimeEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	problemId := int(d.EqualsQuals["problem_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem_timeentry.listProblemTimeEntries", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	entries, _, err := client.Problems.ListTimeEntries(problemId)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem_timeentry.listProblemTimeEntries", "query_error", err)
		return nil, fmt.Errorf("unable to obtain time entries: %v", err)
	}

	for _, entry := range entries.Collection {
		d.StreamListItem(ctx, entry)
	}

	return nil, nil
}
