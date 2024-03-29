package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableRelease() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_release",
		Description: "Obtain information on Releases within the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listReleases,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getRelease,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: releaseColumns(),
	}
}

func releaseColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the release.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "User ID of the agent to whom the release is assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "group_id",
			Description: "ID of the agent group to which the release has been assigned.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "priority",
			Description: "Priority of the release.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "priority_desc",
			Description: "Description of the release priority.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Priority").Transform(releasePriorityDesc),
		},
		{
			Name:        "status",
			Description: "Status of the release.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status_desc",
			Description: "Description of the release status.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Status").Transform(releaseStatusDesc),
		},
		{
			Name:        "release_type",
			Description: "Type of the release.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "release_type_desc",
			Description: "Description of the release type.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("Status").Transform(releaseTypeDesc),
		},
		{
			Name:        "subject",
			Description: "Subject of the release.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "HTML content of the release.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "planned_start_date",
			Description: "Timestamp at which release is planned to start.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "planned_end_date",
			Description: "Timestamp at which release is planned to end.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "work_start_date",
			Description: "Timestamp at which release work started.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "work_end_date",
			Description: "Timestamp at which release work ended.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "department_id",
			Description: "ID of the department initiating the release.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "category",
			Description: "Category of the release.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "sub_category",
			Description: "Sub-category of the release.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "item_category",
			Description: "Item category of the release.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "associated_changes",
			Description: "Array of IDs of Changes associated with the release.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "associated_assets",
			Description: "Array of IDs of Assets associated with the release.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "created_at",
			Description: "Timestamp at which the release was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp at which the release was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getRelease(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.EqualsQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_release.getRelease", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	release, _, err := client.Releases.GetRelease(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_release.getRelease", "query_error", err)
		return nil, fmt.Errorf("unable to obtain release with id %d: %v", id, err)
	}

	return release, nil
}

func listReleases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_release.listReleases", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListReleasesOptions{
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
		releases, res, err := client.Releases.ListReleases(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_release.listReleases", "query_error", err)
			return nil, fmt.Errorf("unable to obtain releases: %v", err)
		}

		for _, release := range releases.Collection {
			d.StreamListItem(ctx, release)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}

// Transform Functions
func releasePriorityDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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

func releaseStatusDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil {
		return "Unknown", nil
	}

	i := input.Value
	switch i.(int) {
	case 1:
		return "Open", nil
	case 2:
		return "On Hold", nil
	case 3:
		return "In Progress", nil
	case 4:
		return "Incomplete", nil
	case 5:
		return "Completed", nil
	default:
		return "Unknown", nil
	}
}

func releaseTypeDesc(_ context.Context, input *transform.TransformData) (interface{}, error) {
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
