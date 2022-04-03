package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableRequester() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_requester",
		Description: "Information about requesters (users) of the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listRequesters,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "email",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getRequester,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: requesterColumns(),
	}
}

func requesterColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "User ID of the requester.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "first_name",
			Description: "First name of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "last_name",
			Description: "Last name of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "job_title",
			Description: "Job title of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "email",
			Description: "Primary email address of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "secondary_emails",
			Description: "Additional/secondary emails associated with the requester.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("AdditionalEmails"),
		},
		{
			Name:        "work_phone_number",
			Description: "Work phone number of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "mobile_phone_number",
			Description: "Mobile phone number of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "department_ids",
			Description: "Array of Unique IDs of the departments associated with the requester.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "active",
			Description: "Indicates if the requester/user is active (enabled)",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "address",
			Description: "Address of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "time_zone",
			Description: "Time zone associated to the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "time_format",
			Description: "Requesters chosen time format (12h or 24h)",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "reporting_manager_id",
			Description: "User ID of the requesters reporting manager.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "language",
			Description: "Language used by the requester. The default language is “en” (English).",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "location_id",
			Description: "Unique ID of the location associated with the requester.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "background_information",
			Description: "Background information of the requester.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "has_logged_in",
			Description: "Set to true if the user has logged in to Freshservice at least once, and false otherwise.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "is_agent",
			Description: "Set to true if the user is also an agent.",
			Type:        proto.ColumnType_BOOL,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the requester was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the requester was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getRequester(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	requester, _, err := client.Requesters.GetRequester(id)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain requester with id %d: %v", id, err)
	}

	return requester, nil
}

func listRequesters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	ia := true
	q := d.KeyColumnQuals
	filter := fs.ListRequestersOptions{
		ListOptions: fs.ListOptions{
			Page:    1,
			PerPage: 30,
		},
		IncludeAgents: &ia,
	}

	if q["email"] != nil {
		e := q["email"].GetStringValue()
		filter.Email = &e
	}

	for {
		users, res, err := client.Requesters.ListRequesters(&filter)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain requesters: %v", err)
		}

		for _, user := range users.Collection {
			d.StreamListItem(ctx, user)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
