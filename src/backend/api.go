package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
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
						with user_data as (select id, username from users where token = $1),
						inserted as (
							insert into rooms (title, token, "creatorId") 
							select $2, $3, id from user_data returning "createdAt", "creatorId"
						)
						select inserted."createdAt", user_data.username 
						from inserted join user_data on inserted."creatorId" = user_data.id
					`
					roomsToken := generateToken()
					row := database.QueryRowContext(context.Background(), query, token, room, roomsToken)
					if row.Err() != nil {
						return Response{Success: false, Data: "such title already exists"}
					}
					var createdAt sql.NullTime
					var username sql.NullString
					row.Scan(&createdAt, &username)
					if !createdAt.Valid || !username.Valid {
						return Response{Success: false, Data: "user with such token is absent"}
					}
					json := map[string]any {
						"title": room,
						"createdAt": createdAt,
						"username": username.String,
						"token": roomsToken,
					}
					return Response{Success: true, Data: json}
				},
			},
			"get": Handler{
				function: func() Response {
					query := `
						select title, rooms."createdAt", username, rooms.token
						from rooms join users on "creatorId" = users.id
					`
					rows, err := database.QueryContext(context.Background(), query)
					if err != nil {
						return Response{Success: false, Data: "inner error"}
					}
					defer rows.Close()
					result := make([]map[string]any, 0)
					for rows.Next() {
						var title, username, token string
						var createdAt time.Time
						err := rows.Scan(&title, &createdAt, &username, &token)
						if err != nil {
							return Response{Success: false, Data: "inner error"}
						}
						json := map[string]any {
							"title": title,
							"createdAt": createdAt,
							"username": username,
							"token": token,
						}
						result = append(result, json)
					}
					return Response{Success: true, Data: result}
				},
			},
		},
	}
}
