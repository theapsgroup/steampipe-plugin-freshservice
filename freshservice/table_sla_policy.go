package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableSlaPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_sla_policy",
		Description: "Obtain information on Service Level Agreement Policies defined in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listSLAs,
		},
		Columns: slaColumns(),
	}
}

func slaColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the SLA Policy.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the SLA Policy.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Short description of the SLA Policy.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "position",
			Description: "Position / Rank of the SLA Policy among other SLA Policies.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "is_default",
			Description: "Set to true if the SLA is the default.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "active",
			Description: "Set to true if the SLA is activated.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "deleted",
			Description: "Set to true if the SLA is deleted.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "sla_targets",
			Description: "Array of SLA Policy targets.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "category",
			Description: "Category the SLA is applicable to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Applicable.Category"),
		},
		{
			Name:        "sub_category",
			Description: "Sub-Category the SLA is applicable to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Applicable.SubCategory"),
		},
		{
			Name:        "item_category",
			Description: "Item Category the SLA is applicable to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Applicable.ItemCategory"),
		},
		{
			Name:        "ticket_type",
			Description: "An array of Ticket Types the SLA is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.TicketType"),
		},
		{
			Name:        "service_items",
			Description: "An array of Service Item IDs the SLA is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.ServiceItems"),
		},
		{
			Name:        "service_categories",
			Description: "An array of Service Category IDs the SLA is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.ServiceCategories"),
		},
		{
			Name:        "department_id",
			Description: "An array of Department IDs the SLA is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.DepartmentIDs"),
		},
		{
			Name:        "group_id",
			Description: "An array of Group IDs the SLA is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.GroupIDs"),
		},
		{
			Name:        "source",
			Description: "An array of Source IDs the SLA is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.Source"),
		},
		{
			Name:        "created_at",
			Description: "SLA Policy creation timestamp.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "SLA Policy updated timestamp.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listSLAs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_sla_policy.listSLAs", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	slas, _, err := client.ServiceLevelAgreements.ListPolicies()
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_sla_policy.listSLAs", "query_error", err)
		return nil, fmt.Errorf("unable to obtain sla policies: %v", err)
	}

	for _, sla := range slas.Collection {
		d.StreamListItem(ctx, sla)
	}

	return nil, nil
}
