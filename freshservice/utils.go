package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"os"
)

func connect(ctx context.Context, d *plugin.QueryData) (*fs.Client, error) {
	domain := os.Getenv("FRESHSERVICE_DOMAIN")
	token := os.Getenv("FRESHSERVICE_TOKEN")

	fsConfig := GetConfig(d.Connection)
	if &fsConfig != nil {
		if fsConfig.Domain != nil {
			domain = *fsConfig.Domain
		}
		if fsConfig.Token != nil {
			token = *fsConfig.Token
		}
	}

	if domain == "" || token == "" {
		errorMsg := ""

		if domain == "" {
			errorMsg += missingConfigOptionError("domain", "FRESHSERVICE_DOMAIN")
		}

		if token == "" {
			errorMsg += missingConfigOptionError("token", "FRESHSERVICE_TOKEN")
		}

		errorMsg += "please set the required values and restart Steampipe"

		return new(fs.Client), fmt.Errorf(errorMsg)
	}

	api, err := fs.NewClient(ctx, domain, token)
	if err != nil {
		return nil, fmt.Errorf("error creating api client for FreshService: %v", err)
	}

	return api, nil
}

// missingConfigOptionError is a utility function for returning parts of error string
func missingConfigOptionError(f string, ev string) string {
	return fmt.Sprintf("configuration option '%s' or Environment Variable '%s' must be set.\n", f, ev)
}
