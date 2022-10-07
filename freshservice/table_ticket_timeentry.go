package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableTicketTimeEntry() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_ticket_timeentry",
		Description: "Obtain time entries for a specific Ticket",
		List: &plugin.ListConfig{
			Hydrate: listTicketTimeEntries,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "ticket_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: ticketTimeEntryColumns(),
	}
}

func ticketTimeEntryColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the time entry.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "start_time",
			Description: "Timestamp when the time entry is added. If a timer, which is in stopped state, is started again, this holds timestamp at which the timer is started again.",
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
			Description: "Set to true if the time entry is billable.",
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
			Name:        "ticket_id",
			Description: "ID of the ticket the time entry belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("ticket_id"),
		},
	}
}

// Hydrate Functions
func listTicketTimeEntries(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	ticketId := int(d.KeyColumnQuals["ticket_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_ticket_timeentry.listTicketTimeEntries", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	entries, _, err := client.Tickets.ListTimeEntries(ticketId)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_ticket_timeentry.listTicketTimeEntries", "query_error", err)
		return nil, fmt.Errorf("unable to obtain time entries: %v", err)
	}

	for _, entry := range entries.Collection {
		d.StreamListItem(ctx, entry)
	}

	return nil, nil
}
