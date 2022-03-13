package freshservice

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-freshservice",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"freshservice_agent":           tableAgent(),
			"freshservice_agent_role":      tableAgentRole(),
			"freshservice_announcement":    tableAnnouncement(),
			"freshservice_asset":           tableAsset(),
			"freshservice_asset_component": tableAssetComponent(),
			"freshservice_asset_contract":  tableAssetContract(),
			"freshservice_asset_type":      tableAssetType(),
		},
	}

	return p
}
