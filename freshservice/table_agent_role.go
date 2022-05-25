package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableAgentRole() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_agent_role",
		Description: "Obtain information about agent roles from the FreshService instance.",
		Columns:     agentRoleColumns(),
		List: &plugin.ListConfig{
			Hydrate: listAgentRoles,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getAgentRole,
			KeyColumns: plugin.SingleColumn("id"),
		},
	}
}

func agentRoleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the role",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the role",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the role",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "default",
			Description: "Set to true if it is a default role, and false otherwise",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the agent role was created",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the agent role was last updated",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

func getAgentRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_agent_role.getAgentRole", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	role, _, err := client.Agents.GetAgentRole(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_agent_role.getAgentRole", "query_error", err)
		return nil, fmt.Errorf("unable to obtain agent role with id %d: %v", id, err)
	}

	return role, nil
}

func listAgentRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_agent_role.listAgentRoles", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListAgentRolesOptions{
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

	for {
		roles, res, err := client.Agents.ListAgentRoles(filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_agent_role.listAgentRoles", "query_error", err)
			return nil, fmt.Errorf("unable to obtain agent roles: %v", err)
		}

		for _, role := range roles.Collection {
			d.StreamListItem(ctx, role)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}
	return nil, nil
}
