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

type ResponsePacket struct {
	Response
	Id   uint64 `json:"id"`
	Type string `json:"type"`
}

func handleRequest(packet RequestFormat, api map[string]map[string]Handler) Response {
	endpoint, exists := api[packet.Service][packet.Method]
	if !exists {
		template := "invalid service and method combination: %s:%s"
		message := fmt.Sprintf(template, packet.Service, packet.Data)
		return Response{Success: false, Data: message}
	}
	arguments, err := CollectArguments(packet.Data, endpoint.fields)
	if err != nil {
		return Response{Success: false, Data: "invalid arguments"}
	}
	result, err := Call(endpoint.function, arguments...)
	if err != nil {
		return Response{Success: false, Data: err.Error()}
	}
	return result.(Response)
}

func onRequest(api map[string]map[string]Handler) func(http.ResponseWriter, *http.Request) {
	connections := make(map[*websocket.Conn]struct{})
	return func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		defer connection.Close()
		connections[connection] = struct{}{}
		for {
			var json map[string]any
			if err := connection.ReadJSON(&json); err != nil {
				message := "Invalid JSON structure"
				err := connection.WriteJSON(Response{Success: false, Data: message})
				if err != nil {
					break
				}
				continue
			}
			packet, err := ParseRequestJson(json)
			if err != nil {
				connection.WriteJSON(Response{Success: false, Data: err.Error()})
				continue
			}
			responsePacket := ResponsePacket{
				Id:       packet.Id,
				Type:     "response",
				Response: handleRequest(packet, api),
			}
			connection.WriteJSON(responsePacket)
			if responsePacket.Success && packet.Service == "rooms" && packet.Method == "create" {
				message := map[string]any{"data": responsePacket.Data}
				for connection := range connections {
					connection.WriteJSON(message)
				}
			}
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
	http.HandleFunc("/", onRequest(apiFactory(database)))
	http.ListenAndServe(address, nil)
}
