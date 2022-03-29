package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTicketTask() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_ticket_task",
		Description: "Obtain tasks based on an associated Ticket",
		List: &plugin.ListConfig{
			Hydrate: listTicketTasks,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "ticket_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: ticketTaskColumns(),
	}
}

func ticketTaskColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the task.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "ID of the agent to whom the task is assigned",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status",
			Description: "Status of the task.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status_desc",
			Description: "Description of the Task Status.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Status").Transform(taskStatusDesc),
		},
		{
			Name:        "due_by",
			Description: "Timestamp that denotes the due date of the task.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "notify_before",
			Description: "Time in seconds before which notification is sent prior to due date.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "title",
			Description: "Title of the task.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the task.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "group_id",
			Description: "Unique ID of the group to which the task is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the task was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the task was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "closed_at",
			Description: "Timestamp at which the task was closed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "ticket_id",
			Description: "Unique ID of the Ticket the task belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("ticket_id"),
		},
	}
}

// Hydrate Functions
func listTicketTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	ticketId := int(d.KeyColumnQuals["ticket_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListTasksOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		tasks, res, err := client.Tickets.ListTasks(ticketId, &filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain tasks: %v", err)
		}

		for _, task := range tasks.Collection {
			d.StreamListItem(ctx, task)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}

// Transform Functions
func taskStatusDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 1:
		return "Open", nil
	case 2:
		return "In Progress", nil
	case 3:
		return "Completed", nil
	default:
		return "Unknown", nil
	}
}
