package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableTicketConversation() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_ticket_conversation",
		Description: "Obtain conversation entries for a specific Ticket",
		List: &plugin.ListConfig{
			Hydrate: listTicketConversations,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "ticket_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: ticketConversationColumns(),
	}
}

func ticketConversationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the conversation.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "user_id",
			Description: "ID of the agent/user who is adding the conversation.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "body",
			Description: "Content of the conversation in HTML.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "body_text",
			Description: "Content of the conversation in plain text.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "source",
			Description: "Denotes the type of the conversation.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "private",
			Description: "Set to true if private.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "incoming",
			Description: "Set to true if a particular conversation should appear as being created from the outside (i.e., not through the web portal).",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "support_email",
			Description: "Email address from which the reply is sent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "ticket_id",
			Description: "Unique ID of the ticket to which this conversation belongs.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "to_emails",
			Description: "Email addresses of agents/users who need to be notified about this conversation.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "attachments",
			Description: "Attachments associated with the conversation..",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Conversation creation timestamp.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Conversation updated timestamp.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listTicketConversations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	ticketId := int(d.KeyColumnQuals["ticket_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListConversationsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		conversations, res, err := client.Tickets.ListConversations(ticketId, &filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain conversations: %v", err)
		}

		for _, conversation := range conversations.Collection {
			d.StreamListItem(ctx, conversation)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
