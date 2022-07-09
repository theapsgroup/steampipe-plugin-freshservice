package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableChange() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_change",
		Description: "Obtain information on Changes raised in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listChanges,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "requester_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getChange,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: changeColumns(),
	}
}

func changeColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "ID of the agent to whom the change is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "description",
			Description: "HTML content of the change.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description_text",
			Description: "Plain text content of the change.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "requester_id",
			Description: "User ID of the initiator/requester of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "group_id",
			Description: "ID of the agent group to which the change is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "department_id",
			Description: "ID of the department initiating the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "priority",
			Description: "Priority of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "priority_desc",
			Description: "Description of the change priority",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Priority").Transform(changePriorityDesc),
		},
		{
			Name:        "status",
			Description: "Status of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status_desc",
			Description: "Description of the change status.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Status").Transform(changeStatusDesc),
		},
		{
			Name:        "impact",
			Description: "Impact of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "impact_desc",
			Description: "Description of the change impact.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Impact").Transform(changeImpactDesc),
		},
		{
			Name:        "risk",
			Description: "Risk of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "risk_desc",
			Description: "Description of the change risk.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Risk").Transform(changeRiskDesc),
		},
		{
			Name:        "type",
			Description: "Type of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "type_desc",
			Description: "Description of the change type.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("ChangeType").Transform(changeTypeDesc),
		},
		{
			Name:        "approval_status",
			Description: "Approval status of the change.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "planned_start_date",
			Description: "Timestamp at which change is starting.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "planned_end_date",
			Description: "Timestamp at which change is ending.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "subject",
			Description: "Subject of the change.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "category",
			Description: "Category of the change.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "sub_category",
			Description: "Sub-category of the change.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "item_category",
			Description: "Item of the change.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the change was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the change was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getChange(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_change.getChange", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	change, _, err := client.Changes.GetChange(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_change.getChange", "query_error", err)
		return nil, fmt.Errorf("unable to obtain change with id %d: %v", id, err)
	}

	return change, nil
}

func listChanges(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_change.listChanges", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	q := d.KeyColumnQuals
	filter := fs.ListChangesOptions{
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

	if q["requester_id"] != nil {
		r := int(q["requester_id"].GetInt64Value())
		filter.RequesterID = &r
	}

	for {
		changes, res, err := client.Changes.ListChanges(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_change.listChanges", "query_error", err)
			return nil, fmt.Errorf("unable to obtain changes: %v", err)
		}

		for _, change := range changes.Collection {
			d.StreamListItem(ctx, change)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}

// Transform Functions
func changePriorityDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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

func changeStatusDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 1:
		return "Open", nil
	case 2:
		return "Planning", nil
	case 3:
		return "Approval", nil
	case 4:
		return "Pending Release", nil
	case 5:
		return "Pending Review", nil
	case 6:
		return "Closed", nil
	default:
		return "Unknown", nil
	}
}

func changeTypeDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 1:
		return "Minor", nil
	case 2:
		return "Standard", nil
	case 3:
		return "Major", nil
	case 4:
		return "Emergency", nil
	default:
		return "Unknown", nil
	}
}

func changeImpactDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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
	default:
		return "Unknown", nil
	}
}

func changeRiskDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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
		return "Very High", nil
	default:
		return "Unknown", nil
	}
}
