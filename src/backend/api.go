package main

import (
	"context"
	"database/sql"
)

type Handler struct {
	fields   []string
	function any
}

type Response struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func apiFactory(database *sql.DB) map[string]map[string]Handler {
	return map[string]map[string]Handler{
		"users": {
			"exists": Handler{
				function: func(token string) Response {
					var exists bool
					query := "select exists(select 1 from users where token = $1)"
					row := database.QueryRowContext(context.Background(), query, token)
					err := row.Scan(&exists)
					if err != nil {
						return Response{Success: false, Data: "Inner error"}
					}
					return Response{Success: true, Data: exists}
				},
			},
		},
	}
}
