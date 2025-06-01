package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
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

func call(function any, args ...any) error {
	value := reflect.ValueOf(function)
	if value.Kind() != reflect.Func {
		return errors.New("value is not a function")
	}
	if value.Type().NumIn() != len(args) {
		return errors.New("invalid arguments count")
	}
	arguments := make([]reflect.Value, value.Type().NumIn())
	for index := range args {
		functionArgumentType := value.Type().In(index)
		argumentValue := reflect.ValueOf(args[index])
		if !argumentValue.Type().ConvertibleTo(functionArgumentType) {
			return errors.New("invalid arguments format")
		}
		converted := argumentValue.Convert(functionArgumentType)
		arguments[index] = converted
	}
	value.Call(arguments)
	return nil
}

func readFields(json map[string]any, fields []string) (map[string]any, error) {
	response := make(map[string]any)
	for _, field := range fields {
		value, ok := json[field]
		if !ok {
			return nil, fmt.Errorf("key '%s' is absent", field)
		}
		response[field] = value
	}
	return response, nil
}

type RequestFormat struct {
	service, method string
	data any
}

func parseRequestJson(json map[string]any) (RequestFormat, error) {
	packet, err := readFields(json, []string{"service", "method"})
	if err != nil {
		return RequestFormat{}, err
	}
	service, ok := packet["service"].(string)
	if !ok {
		return RequestFormat{}, errors.New("invalid service key")
	}
	method, ok := packet["method"].(string)
	if !ok {
		return RequestFormat{}, errors.New("invalid method key")
	}
	return RequestFormat{service: service, method: method, data: json["data"]}, nil
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
			connection.ReadJSON(&json)
			packet, err := parseRequestJson(json)
			if err != nil {
				connection.WriteJSON(Response{ Success: false, Data: err.Error() })
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
