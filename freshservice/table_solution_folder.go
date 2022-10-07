package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableSolutionFolder() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_solution_folder",
		Description: "Obtain Solution Folders from the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listSolutionFolders,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: solutionFolderColumns(),
	}
}

func solutionFolderColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the solution folder.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "name",
			Description: "Name of the solution folder.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description about the solution folder.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "position",
			Description: "The rank of the solution folder in the folder listing.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "default_category",
			Description: "Set to true if the solution folder is the default one.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "category_id",
			Description: "ID of the category under which the solution folder is listed.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "visibility",
			Description: "Accessibility of this solution folder.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "department_ids",
			Description: "Array of IDs of the departments to which this solution folder is visible.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "group_ids",
			Description: "Array of IDs of the agent groups to which this solution folder is visible.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "requester_group_ids",
			Description: "Array of IDs of requester groups to which this solution folder is visible.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "manage_by_group_ids",
			Description: "Array of IDs of groups which manage this solution folder.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the solution folder was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the solution folder was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func listSolutionFolders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_solution_folder.listSolutionFolders", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListSolutionFoldersOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	q := d.KeyColumnQuals

	if q["id"] != nil {
		folderId := int(q["id"].GetInt64Value())

		folder, _, err := client.Solutions.GetSolutionFolder(folderId)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_solution_folder.listSolutionFolders", "query_error", err)
			return nil, fmt.Errorf("unable to obtain solution folder with id %d: %v", folderId, err)
		}

		d.StreamListItem(ctx, folder)
	} else {
		for {
			folders, res, err := client.Solutions.ListSolutionFolders(&filter)
			if err != nil {
				plugin.Logger(ctx).Error("freshservice_solution_folder.listSolutionFolders", "query_error", err)
				return nil, fmt.Errorf("unable to obtain solution folders: %v", err)
			}

			for _, folder := range folders.Collection {
				d.StreamListItem(ctx, folder)
			}

			if res.Header.Get("link") == "" {
				break
			}

			filter.Page += 1
		}
	}

	return nil, nil
}
