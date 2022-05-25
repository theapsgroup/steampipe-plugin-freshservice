package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableTicket() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_ticket",
		Description: "Obtain information on Tickets raised in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listTickets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "email",
					Require: plugin.Optional,
				},
				{
					Name:    "requester_id",
					Require: plugin.Optional,
				},
				{
					Name:    "type",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getTicket,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: ticketColumns(),
	}
}

func ticketColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the ticket.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "subject",
			Description: "Subject of the ticket.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "HTML content of the ticket.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description_text",
			Description: "Plain text content of the ticket.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "requester_id",
			Description: "User ID of the requester.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "email",
			Description: "Email address of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "phone",
			Description: "Phone number of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "status",
			Description: "Status of the ticket.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status_desc",
			Description: "Description of the Ticket Status.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Status").Transform(ticketStatusDesc),
		},
		{
			Name:        "priority",
			Description: "Priority of the ticket.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "priority_desc",
			Description: "Description of the Ticket Priority",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Priority").Transform(ticketPriorityDesc),
		},
		{
			Name:        "category",
			Description: "Ticket Category.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "sub_category",
			Description: "Ticket sub category.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "item_category",
			Description: "Ticket item category.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "type",
			Description: "Helps categorize the ticket according to the different kinds of issues your support team deals with.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "urgency",
			Description: "Ticket urgency.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "impact",
			Description: "Ticket impact.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "responder_id",
			Description: "ID of the agent to whom the ticket has been assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "fr_due_by",
			Description: "Timestamp that denotes when the first response is due.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "fr_escalated",
			Description: "Set to true if the ticket has been escalated as a result of the first response time being breached.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "due_by",
			Description: "Timestamp that denotes when the ticket is due to be resolved.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "is_escalated",
			Description: "Set to true if the ticket has been escalated for any reason.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "deleted",
			Description: "Set to true if the ticket has been deleted/trashed.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "department_id",
			Description: "ID of the department to which this ticket belongs.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "group_id",
			Description: "ID of the group to which the ticket has been assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "spam",
			Description: "Set to true if the ticket has been marked as spam.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "source",
			Description: "The channel through which the ticket was created.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "tags",
			Description: "Tags that have been associated with the ticket.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "attachments",
			Description: "Ticket attachments.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Ticket creation timestamp.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Ticket updated timestamp.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getTicket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_ticket.getTicket", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	ticket, _, err := client.Tickets.GetTicket(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_ticket.getTicket", "query_error", err)
		return nil, fmt.Errorf("unable to obtain ticket with id %d: %v", id, err)
	}

	return ticket, nil
}

func listTickets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_ticket.listTickets", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListTicketsOptions{
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

	if q["email"] != nil {
		e := q["email"].GetStringValue()
		filter.Email = &e
	}

	if q["requester_id"] != nil {
		r := int(q["requester_id"].GetInt64Value())
		filter.RequesterID = &r
	}

	if q["type"] != nil {
		t := q["type"].GetStringValue()
		filter.Type = &t
	}

	for {
		tickets, res, err := client.Tickets.ListTickets(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_ticket.listTickets", "query_error", err)
			return nil, fmt.Errorf("unable to obtain tickets: %v", err)
		}

		for _, ticket := range tickets.Collection {
			d.StreamListItem(ctx, ticket)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}

// Transform Functions
func ticketStatusDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 2:
		return "Open", nil
	case 3:
		return "Pending", nil
	case 4:
		return "Resolved", nil
	case 5:
		return "Closed", nil
	default:
		return "Unknown", nil
	}
}

func ticketPriorityDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 1:
		return "Low", nil
	case 2:
		return "Medium", nil
	case 3:
		return "High", nil
	case 4:
		return "Urgent", nil
	default:
		return "Unknown", nil
	}
}
