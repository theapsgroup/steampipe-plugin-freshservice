package freshservice

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-freshservice",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"freshservice_agent":               tableAgent(),
			"freshservice_agent_role":          tableAgentRole(),
			"freshservice_announcement":        tableAnnouncement(),
			"freshservice_asset":               tableAsset(),
			"freshservice_asset_component":     tableAssetComponent(),
			"freshservice_asset_contract":      tableAssetContract(),
			"freshservice_asset_type":          tableAssetType(),
			"freshservice_business_hour":       tableBusinessHour(),
			"freshservice_change":              tableChange(),
			"freshservice_change_note":         tableChangeNote(),
			"freshservice_department":          tableDepartment(),
			"freshservice_location":            tableLocation(),
			"freshservice_requester":           tableRequester(),
			"freshservice_product":             tableProduct(),
			"freshservice_ticket":              tableTicket(),
			"freshservice_ticket_conversation": tableTicketConversation(),
			"freshservice_ticket_task":         tableTicketTask(),
			"freshservice_ticket_timeentry":    tableTicketTimeEntry(),
		},
	}

	return p
}
