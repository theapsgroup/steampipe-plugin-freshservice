package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableProblemTask() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_problem_task",
		Description: "Obtain tasks based on an associated Problem",
		List: &plugin.ListConfig{
			Hydrate: listProblemTasks,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "problem_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: problemTaskColumns(),
	}
}

func problemTaskColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the task.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "ID of the agent to whom the task is assigned",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status",
			Description: "Status of the task.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status_desc",
			Description: "Description of the Task Status.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Status").Transform(taskStatusDesc),
		},
		{
			Name:        "due_date",
			Description: "Timestamp that denotes the due date of the task.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "notify_before",
			Description: "Time in seconds before which notification is sent prior to due date.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "title",
			Description: "Title of the task.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the task.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "group_id",
			Description: "Unique ID of the group to which the task is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the task was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the task was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "closed_at",
			Description: "Timestamp at which the task was closed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "problem_id",
			Description: "Unique ID of the Problem the task belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("problem_id"),
		},
	}
}

// Hydrate Functions
func listProblemTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	problemId := int(d.KeyColumnQuals["problem_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListTasksOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	for {
		tasks, res, err := client.Problems.ListTasks(problemId, &filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain tasks: %v", err)
		}

		for _, task := range tasks.Collection {
			d.StreamListItem(ctx, task)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
