package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const host = "0.0.0.0"
const port = "8080"

var address = net.JoinHostPort(host, port)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1 << 10,
	WriteBufferSize:   1 << 10,
	HandshakeTimeout:  time.Second * 3,
	EnableCompression: true,
	CheckOrigin:       func(*http.Request) bool { return true },
}

var url = "postgres://postgres:marcus@localhost:5432/global_chat"

func mainPageRequest(_ map[string]map[string]Handler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		defer connection.Close()
		for {
			var json map[string]any
			if err := connection.ReadJSON(&json); err != nil {
				message := "Invalid JSON structure"
				connection.WriteJSON(Response{Success: false, Data: message})
			}
			packet, err := parseRequestJson(json)
			if err != nil {
				connection.WriteJSON(Response{Success: false, Data: err.Error()})
				continue
			}
			fmt.Println(packet)
		}
	}
}

func main() {
	database, err := sql.Open("pgx", url)
	if err != nil {
		fmt.Fprint(os.Stdout, "Cannot initialise database")
	}
	http.HandleFunc("/", mainPageRequest(apiFactory(database)))
	http.ListenAndServe(address, nil)
}
