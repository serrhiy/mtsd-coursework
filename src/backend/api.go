package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
)

type Handler struct {
	fields   []string
	function any
}

type Response struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

var keyLength = 16

func generateToken() string {
	buffer := make([]byte, keyLength)
	rand.Read(buffer)
	return hex.EncodeToString(buffer)
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
						return Response{Success: false, Data: "inner error"}
					}
					return Response{Success: true, Data: exists}
				},
			},
			"create": Handler{
				function: func(username string) Response {
					query := "insert into users (token, username) values ($1, $2)"
					token := generateToken()
					_, err := database.ExecContext(context.Background(), query, token, username)
					if err != nil {
						return Response{Success: false, Data: "inner error"}
					}
					return Response{Success: true, Data: token}
				},
			},
			"update": Handler{
				fields: []string{"token", "username"},
				function: func(token, username string) Response {
					query := "update users set username = $1 where token = $2"
					result, err := database.ExecContext(context.Background(), query, username, token)
					if n, err := result.RowsAffected(); err == nil && n == 0 {
						return Response{Success: false, Data: "user with such token is absent"}
					}
					if err != nil {
						return Response{Success: false, Data: "inner error"}
					}
					return Response{Success: true}
				},
			},
		},
		"rooms": {
			"create": Handler{
				fields: []string{"token", "room"},
				function: func(token, room string) Response {
					query := `
						INSERT INTO rooms (title, "creatorId")
						SELECT $1, id FROM users WHERE token = $2;
					`
					result, err := database.ExecContext(context.Background(), query, room, token)
					if err != nil {
						return Response{Success: false, Data: "such title already exists"}
					}
					if n, err := result.RowsAffected(); err == nil && n == 0 {
						return Response{Success: false, Data: "user with such token is absent"}
					}
					return Response{Success: true}
				},
			},
		},
	}
}
