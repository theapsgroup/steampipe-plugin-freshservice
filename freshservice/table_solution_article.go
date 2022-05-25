package freshservice

import (
	"context"
	"fmt"
	fs "github.com/theapsgroup/go-freshservice/freshservice"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableSolutionArticle() *plugin.Table {
	return &plugin.Table{
		Name:        "freshservice_solution_article",
		Description: "Obtain information on Solution Articles stored in the FreshService instance.",
		List: &plugin.ListConfig{
			Hydrate: listSolutionArticles,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "folder_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getSolutionArticle,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: solutionArticleColumns(),
	}
}

func solutionArticleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Description: "Unique ID of the solution article.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "title",
			Description: "Title of the solution article.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "description",
			Description: "Description of the solution article.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "position",
			Description: "The rank of the solution article in the article listing.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "article_type",
			Description: "The type of the article (1 permanent, 2 work-around).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "folder_id",
			Description: "ID of the folder under which the article is listed.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "category_id",
			Description: "ID of the category under which the solution article is belongs.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "status",
			Description: "Status of the solution article (1 draft, 2 published).",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "approval_status",
			Description: "Approval status of the solution article.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "thumbs_up",
			Description: "Number of up-votes for the solution article.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "thumbs_down",
			Description: "Number of down-votes for the solution article.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "agent_id",
			Description: "ID of the user who created the solution article.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "views",
			Description: "Total hit count of the page visits for this solution article.",
			Type:        proto.ColumnType_INT,
		},
		{
			Name:        "tags",
			Description: "Array of tags that have been associated with the solution article.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "keywords",
			Description: "Array of keywords that have been associated with the solution article.",
			Type:        proto.ColumnType_JSON,
		},
		{
			Name:        "url",
			Description: "External url for the Solution Article.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "review_date",
			Description: "Timestamp at which the solution article needs to be reviewed.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "created_at",
			Description: "Timestamp when the solution article was created.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "updated_at",
			Description: "Timestamp when the solution article was last updated.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
	}
}

// Hydrate Functions
func getSolutionArticle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := int(d.KeyColumnQuals["id"].GetInt64Value())

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_solution_article.getSolutionArticle", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	article, _, err := client.Solutions.GetSolutionArticle(id)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_solution_article.getSolutionArticle", "query_error", err)
		return nil, fmt.Errorf("unable to obtain solution article with id %d: %v", id, err)
	}

	return article, nil
}

func listSolutionArticles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("freshservice_solution_article.listSolutionArticles", "connection_error", err)
		return nil, fmt.Errorf("unable to create FreshService client: %v", err)
	}

	filter := fs.ListSolutionArticlesOptions{
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
	if q["folder_id"] != nil {
		fid := int(q["folder_id"].GetInt64Value())
		filter.FolderID = fid
	}

	for {
		articles, res, err := client.Solutions.ListSolutionArticles(&filter)
		if err != nil {
			plugin.Logger(ctx).Error("freshservice_solution_article.listSolutionArticles", "query_error", err)
			return nil, fmt.Errorf("unable to obtain solution articles: %v", err)
		}

		for _, article := range articles.Collection {
			d.StreamListItem(ctx, article)
		}

		if res.Header.Get("link") == "" {
			break
		}

		filter.Page += 1
	}

	return nil, nil
}
