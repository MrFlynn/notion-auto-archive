package main

import (
	"context"

	"github.com/dstotijn/go-notion"
)

// NotionHandler is container that holds a notion client and other associated information.
type NotionHandler struct {
	client *notion.Client
}

// NewNotionHandler creates a new NotionHandler struct with an initialized handler.
func NewNotionHandler(key string) (handler *NotionHandler) {
	handler = &NotionHandler{
		client: notion.NewClient(key),
	}

	return
}

// GetTasksOnBoard returns a slice ot page (task) objects on that match the Selectors.SourceColumn attribute of the
// TaskBoard parameter.
func (n *NotionHandler) GetTasksOnBoard(ctx context.Context, board TaskBoard) (tasks []notion.Page, err error) {
	var (
		results notion.DatabaseQueryResponse

		query = &notion.DatabaseQuery{
			Filter: &notion.DatabaseQueryFilter{
				Property: board.Selectors.ColumnName,
				Select: &notion.SelectDatabaseQueryFilter{
					Equals: board.Selectors.SourceColumn,
				},
			},
		}
	)

	tasks = make([]notion.Page, 0, 10)

	for {
		results, err = n.client.QueryDatabase(ctx, board.ID, query)
		if err != nil {
			return
		}

		tasks = append(tasks, results.Results...)

		if !results.HasMore {
			break
		}

		query.StartCursor = *results.NextCursor
	}

	return
}

// MoveTasks moves all the individual pages (tasks) to the column specified in Selectors.TargetColumn.
func (n *NotionHandler) MoveTasks(ctx context.Context, board TaskBoard, tasks ...notion.Page) (err error) {
	for _, task := range tasks {
		params := notion.UpdatePageParams{
			DatabasePageProperties: &notion.DatabasePageProperties{
				board.Selectors.ColumnName: notion.DatabasePageProperty{
					Select: &notion.SelectOptions{
						Name: board.Selectors.TargetColumn,
					},
				},
			},
		}

		_, err = n.client.UpdatePageProps(ctx, task.ID, params)
		if err != nil {
			return
		}
	}

	return
}
