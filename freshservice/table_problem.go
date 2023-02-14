package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableProblem() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_problem",
		Description: "Obtain information on Problems raised in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listProblems,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getProblem,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: problemColumns(),
	}
}

func problemColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "User ID of the agent to whom the problem is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "requester_id",
			Description: "User ID of the requester.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "group_id",
			Description: "ID of the agent group to which the problem has been assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "description",
			Description: "HTML content of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description_text",
			Description: "Plain text content of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "priority",
			Description: "Priority of the problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "priority_desc",
			Description: "Description of the problems priority",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Priority").Transform(problemPriorityDesc),
		},
		{
			Name:        "status",
			Description: "Status of the problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status_desc",
			Description: "Description of the problems status.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Status").Transform(problemStatusDesc),
		},
		{
			Name:        "impact",
			Description: "Impact of the problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "impact_desc",
			Description: "Description of the problems impact.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Impact").Transform(problemImpactDesc),
		},
		{
			Name:        "known_error",
			Description: "Set to true if the problem is a known issue/problem/error.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "subject",
			Description: "Subject of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "due_by",
			Description: "Timestamp at which problem is due to be resolved by.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "department_id",
			Description: "ID of the department initiating the problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "category",
			Description: "Category of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "sub_category",
			Description: "Sub-category of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "item_category",
			Description: "Item category of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "associated_change",
			Description: "ID of the change associated with the problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "problem_cause",
			Description: "Cause of the problem.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("AnalysisFields.ProblemCause.DescriptionText"),
		},
		{
			Name:        "problem_symptom",
			Description: "Symptom(s) of the problem.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("AnalysisFields.ProblemSymptom.DescriptionText"),
		},
		{
			Name:        "problem_impact",
			Description: "Impact of the problem (textual description).",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("AnalysisFields.ProblemImpact.DescriptionText"),
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the problem was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the problem was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getProblem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.EqualsQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem.getProblem", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	problem, _, err := client.Problems.GetProblem(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem.getProblem", "query_error", err)
		return nil, fmt.Errorf("unable to obtain problem with id %d: %v", id, err)
	}

	return problem, nil
}

func listProblems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_problem.listProblems", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListProblemsOptions{
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
		problems, res, err := client.Problems.ListProblems(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_problem.listProblems", "query_error", err)
			return nil, fmt.Errorf("unable to obtain releases: %v", err)
		}

		for _, problem := range problems.Collection {
			d.StreamListItem(ctx, problem)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}

// Transform Functions
func problemPriorityDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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

func problemStatusDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 1:
		return "Open", nil
	case 2:
		return "Change Requested", nil
	case 3:
		return "Closed", nil
	default:
		return "Unknown", nil
	}
}

func problemImpactDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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
