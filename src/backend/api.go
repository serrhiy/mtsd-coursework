package main

import (
	"database/sql"
	"fmt"
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
		"user": {
			"exists": Handler{
				fields: []string{"token"},
				function: func(token string) Response {
					fmt.Println(token)
					return Response{Success: true, Data: false}
				},
			},
		},
	}
}
