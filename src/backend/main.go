package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

func onRequest(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err, request.Header["Origin"])
		return
	}
	defer connection.Close()
	connection.WriteMessage(websocket.TextMessage, []byte("Hello world!"))
}

func main() {
	http.HandleFunc("/", onRequest)
	http.ListenAndServe(address, nil)
}
