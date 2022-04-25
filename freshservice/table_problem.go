package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			Description: "Unique ID of the Problem.",
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
			Description: "Description of the Problem Priority",
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
			Description: "Description of the Problem Status.",
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
			Description: "Description of the Problem Impact.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Impact").Transform(problemImpactDesc),
		},
		{
			Name:        "known_error",
			Description: "States that the problem is known issue or not. true or false.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "subject",
			Description: "Subject of the problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "due_by",
			Description: "Timestamp at which Problem due ends.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "department_id",
			Description: "ID of the department initiating the Problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "category",
			Description: "Category of the Problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "sub_category",
			Description: "Sub-category of the Problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "item_category",
			Description: "Item of the Problem.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "associated_change",
			Description: "ID of the Change associated with the Problem.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "problem_cause",
			Description: "Cause of the Problem.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("AnalysisFields.ProblemCause.DescriptionText"),
		},
		{
			Name:        "problem_symptom",
			Description: "Symptom of the Problem.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("AnalysisFields.ProblemSymptom.DescriptionText"),
		},
		{
			Name:        "problem_impact",
			Description: "Impact of the Problem (textual description).",
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
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	problem, _, err := client.Problems.GetProblem(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain problem with id %d: %v", id, err)
	}

	return problem, nil
}

func listProblems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListProblemsOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		problems, res, err := client.Problems.ListProblems(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain tickets: %v", err)
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
