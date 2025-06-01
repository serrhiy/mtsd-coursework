package main

import (
	"database/sql"
	"errors"
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

func CollectArguments(data any, fields []string) ([]any, error) {
	if fields == nil {
		return []any{data}, nil
	}
	object, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("invalid arguments")
	}
	result := make([]any, len(object))
	for index, field := range fields {
		value, exists := object[field]
		if !exists {
			return nil, errors.New("key " + field + " is absent")
		}
		result[index] = value
	}
	return result, nil
}

func mainPageRequest(api map[string]map[string]Handler) func(http.ResponseWriter, *http.Request) {
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
			responsePacket := ResponsePacket{Id: packet.id, Type: "response"}
			endpoint, exists := api[packet.service][packet.method]
			if !exists {
				template := "invalid service and method combination: %s:%s"
				message := fmt.Sprintf(template, packet.service, packet.data)
				responsePacket.Response = Response{Success: false, Data: message}
				connection.WriteJSON(responsePacket)
				continue
			}
			arguments, err := CollectArguments(packet.data, endpoint.fields)
			if err != nil {
				responsePacket.Response = Response{Success: false, Data: "invalid arguments"}
				connection.WriteJSON(responsePacket)
				continue
			}
			result, err := Call(endpoint.function, arguments...)
			if err != nil {
				responsePacket.Response = Response{Success: false, Data: "invalid arguments"}
				connection.WriteJSON(responsePacket)
				continue
			}
			responsePacket.Response = result.(Response)
			connection.WriteJSON(responsePacket)
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
