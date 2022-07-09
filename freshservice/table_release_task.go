package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableReleaseTask() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_release_task",
		Description: "Obtain tasks based on an associated Release",
		List: &plugin.ListConfig{
			Hydrate: listReleaseTasks,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "release_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: releaseTaskColumns(),
	}
}

func releaseTaskColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the task.",
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
			Description: "Description of the task status.",
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
			Description: "ID of the group to which the task is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the task was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the task was updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "closed_at",
			Description: "Timestamp when the task was closed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "release_id",
			Description: "ID of the release the task belongs to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("release_id"),
		},
	}
}

// Hydrate Functions
func listReleaseTasks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	releaseId := int(d.KeyColumnQuals["release_id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_release_task.listReleaseTasks", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListTasksOptions{
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
		tasks, res, err := client.Releases.ListTasks(releaseId, &filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_release_task.listReleaseTasks", "query_error", err)
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
