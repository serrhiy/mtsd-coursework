package main

import (
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
					return Response{Success: true, Data: false}
				},
			},
		},
	}
}
