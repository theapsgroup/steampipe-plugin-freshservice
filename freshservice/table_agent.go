package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableAgent() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_agent",
		Description: "Information about agents (operators) of the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listAgents,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "email",
					Require: plugin.Optional,
				},
				{
					Name:    "active",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAgent,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: agentColumns(),
	}
}

func agentColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "User ID of the agent.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "first_name",
			Description: "First name of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "last_name",
			Description: "Last name of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "occasional",
			Description: "True if the agent is an occasional agent, and false if full-time agent.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "active",
			Description: "Indicates if the agent is active (enabled)",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "job_title",
			Description: "Job title of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "email",
			Description: "Email address of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "work_phone_number",
			Description: "Work phone number of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "mobile_phone_number",
			Description: "Mobile phone number of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "reporting_manager_id",
			Description: "User ID of the agent's reporting manager.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "address",
			Description: "Address of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "time_zone",
			Description: "Time zone associated to the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "time_format",
			Description: "Agents chosen time format (12h or 24h)",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "language",
			Description: "Language used by the agent. The default language is “en” (English).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "location_id",
			Description: "Unique ID of the location associated with the agent.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "background_information",
			Description: "Background information of the agent.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "scoreboard_level_id",
			Description: "Unique ID of the level of the agent in the Arcade. 1 (Beginner), 2 (Intermediate), 3 (Professional), 4 (Expert), 5 (Master), 6 (Guru).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "member_of",
			Description: "Unique IDs of the groups that the agent is a member of.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "roles",
			Description: "Array of roles assigned to the agent, each individual role is a hash that contains the attributes role_id, assignment_scope & groups.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "department_ids",
			Description: "Array of Unique IDs of the departments associated with the agent",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "last_login_at",
			Description: "Timestamp of the agent's last successful login.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "last_active_at",
			Description: "Timestamp of the agent's recent activity.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "has_logged_in",
			Description: "Set to true if the user has logged in to Freshservice at least once, and false otherwise.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the agent was created",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the agent was last updated",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getAgent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService Agent client: %v", err)
	}

	agent, _, err := client.Agents.GetAgent(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain agent with id %d: %v", id, err)
	}

	return agent, nil
}

func listAgents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService Agent client: %v", err)
	}

	filter := fs.ListAgentsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}
	q := d.KeyColumnQuals

	if q["email"] != nil {
		e := q["email"].GetStringValue()
		filter.Email = &e
	}

	if q["active"] != nil {
		a := q["active"].GetBoolValue()
		filter.Active = &a
	}

	for {
		agents, res, err := client.Agents.ListAgents(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain agents: %v", err)
		}

		for _, agent := range agents.Collection {
			d.StreamListItem(ctx, agent)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}
	return nil, nil
}
