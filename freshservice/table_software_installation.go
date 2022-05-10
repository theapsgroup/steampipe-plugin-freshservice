package freshservice

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableSoftwareInstallation() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_software_installation",
		Description: "Obtain information about Installations of Software registered in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listSoftwareInstallations,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "software_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: softwareInstallationColumns(),
	}
}

func softwareInstallationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the software installation.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "installation_machine_id",
			Description: "Display ID of device (asset) where the software is installed.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "installation_path",
			Description: "Path where the software is installed.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "version",
			Description: "Version of the installed software.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "user_id",
			Description: "ID of the user using the device.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "department_id",
			Description: "ID of the department the device belongs to.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "installation_date",
			Description: "Date and time when the software was installed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "created_at",
			Description: "Date and time when the installation was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Date and time when the installation was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "software_id",
			Description: "Unique ID of the software this installation belong to.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromQual("software_id"),
		},
	}
}

// Hydrate Functions
func listSoftwareInstallations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	q := d.KeyColumnQuals
	s := int(q["software_id"].GetInt64Value())

	installs, _, err := client.Software.ListInstallations(s)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain software installations: %v", err)
	}

	for _, install := range installs.Collection {
		d.StreamListItem(ctx, install)
	}

	return nil, nil
}
