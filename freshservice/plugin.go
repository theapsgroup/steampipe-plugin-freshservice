package freshservice

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			"freshservice_agent":                 tableAgent(),
			"freshservice_agent_role":            tableAgentRole(),
			"freshservice_announcement":          tableAnnouncement(),
			"freshservice_asset":                 tableAsset(),
			"freshservice_asset_component":       tableAssetComponent(),
			"freshservice_asset_contract":        tableAssetContract(),
			"freshservice_asset_type":            tableAssetType(),
			"freshservice_business_hour":         tableBusinessHour(),
			"freshservice_change":                tableChange(),
			"freshservice_change_note":           tableChangeNote(),
			"freshservice_contract":              tableContract(),
			"freshservice_contract_type":         tableContractType(),
			"freshservice_department":            tableDepartment(),
			"freshservice_location":              tableLocation(),
			"freshservice_requester":             tableRequester(),
			"freshservice_problem":               tableProblem(),
			"freshservice_problem_note":          tableProblemNote(),
			"freshservice_problem_task":          tableProblemTask(),
			"freshservice_problem_timeentry":     tableProblemTimeEntry(),
			"freshservice_product":               tableProduct(),
			"freshservice_purchase_order":        tablePurchaseOrder(),
			"freshservice_release":               tableRelease(),
			"freshservice_release_note":          tableReleaseNote(),
			"freshservice_release_task":          tableReleaseTask(),
			"freshservice_release_timeentry":     tableReleaseTimeEntry(),
			"freshservice_service":               tableService(),
			"freshservice_sla_policy":            tableSlaPolicy(),
			"freshservice_software":              tableSoftware(),
			"freshservice_software_installation": tableSoftwareInstallation(),
			"freshservice_software_user":         tableSoftwareUser(),
			"freshservice_solution_article":      tableSolutionArticle(),
			"freshservice_solution_category":     tableSolutionCategory(),
			"freshservice_solution_folder":       tableSolutionFolder(),
			"freshservice_ticket":                tableTicket(),
			"freshservice_ticket_conversation":   tableTicketConversation(),
			"freshservice_ticket_task":           tableTicketTask(),
			"freshservice_ticket_timeentry":      tableTicketTimeEntry(),
			"freshservice_vendor":                tableVendor(),
		},
	}

	return p
}
