package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableSoftwareUser() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_software_user",
		Description: "Obtain information Users assigned to Software.",
		List: &plugin.ListConfig{
			Hydrate: listSoftwareUsers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "software_id",
					Require: plugin.Required,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: softwareUserColumns(),
	}
}

func softwareUserColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "ID of the software-user combination.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "user_id",
			Description: "ID of the user using the software.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "license_id",
			Description: "Display ID of the allocated software license contract.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "allocated_date",
			Description: "Timestamp when the license was allocated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "first_used",
			Description: "Timestamp when the software was first used by the user.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "last_used",
			Description: "Timestamp when the software was last used by the user.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the installation was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the installation was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "software_id",
			Description: "ID of the software this installation belong to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("software_id"),
		},
	}
}

// Hydrate Functions
func listSoftwareUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_software_user.listSoftwareUsers", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListSoftwareUsersOptions{
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

	q := d.KeyColumnQuals
	s := int(q["software_id"].GetInt64Value())

	if q["id"] != nil {
		u := int(q["id"].GetInt64Value())
		user, _, err := client.Software.GetSoftwareUser(s, u)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_software_user.listSoftwareUsers", "query_error", err)
			return nil, fmt.Errorf("unable to obtain software user: %v", err)
		}

		d.StreamListItem(ctx, user)
	} else {
		for {
			users, res, err := client.Software.ListSoftwareUsers(s, &filter)
			if err != nil {
				plugin.Logger(ctx).Error("freshservice_software_user.listSoftwareUsers", "query_error", err)
				return nil, fmt.Errorf("unable to obtain software users: %v", err)
			}

			for _, user := range users.Collection {
				d.StreamListItem(ctx, user)
			}

			if res.Header.Get("link") == "" {
				break
			}

			filter.Page += 1
		}
	}

	return nil, nil
}
