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

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1 << 10,
	WriteBufferSize:   1 << 10,
	HandshakeTimeout:  time.Second * 3,
	EnableCompression: true,
	CheckOrigin:       func(*http.Request) bool { return true },
}

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
			packet, err := ParseRequestJson(json)
			if err != nil {
				connection.WriteJSON(Response{Success: false, Data: err.Error()})
				continue
			}
			fmt.Println(packet)
		}
	}
}

func main() {
	config, err := GetConfig("config.json")
	if err != nil {
		fmt.Fprint(os.Stdout, "Cannot read config file\n")
	}
	database, err := sql.Open("pgx", DatabaseURL(config.Database))
	if err != nil {
		fmt.Fprint(os.Stdout, "Cannot initialise database\n")
	}
	address := net.JoinHostPort(config.Network.Host, config.Network.Port)
	http.HandleFunc("/", mainPageRequest(apiFactory(database)))
	http.ListenAndServe(address, nil)
}
