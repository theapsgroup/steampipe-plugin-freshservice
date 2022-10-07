package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
			Description: "ID of the policy.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the policy.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Short description of the policy.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "position",
			Description: "Position (ranking) of the policy among other policies.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "is_default",
			Description: "Set to true if the policy is the default.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "active",
			Description: "Set to true if the policy is activated.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "deleted",
			Description: "Set to true if the policy is deleted.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "sla_targets",
			Description: "Array of policy targets.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "category",
			Description: "Category the policy is applicable to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Applicable.Category"),
		},
		{
			Name:        "sub_category",
			Description: "Sub-Category the policy is applicable to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Applicable.SubCategory"),
		},
		{
			Name:        "item_category",
			Description: "Item Category the policy is applicable to.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Applicable.ItemCategory"),
		},
		{
			Name:        "ticket_type",
			Description: "Array of ticket types the policy is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.TicketType"),
		},
		{
			Name:        "service_items",
			Description: "Array of service item IDs the policy is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.ServiceItems"),
		},
		{
			Name:        "service_categories",
			Description: "Array of service category IDs the policy is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.ServiceCategories"),
		},
		{
			Name:        "department_id",
			Description: "Array of department IDs the policy is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.DepartmentIDs"),
		},
		{
			Name:        "group_id",
			Description: "Array of group IDs the policy is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.GroupIDs"),
		},
		{
			Name:        "source",
			Description: "Array of source IDs the policy is applicable to.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("Applicable.Source"),
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the policy was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the policy was last updated.",
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
