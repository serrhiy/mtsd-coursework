package main

import (
	"io"
	"net"
	"net/http"
)

const host = "0.0.0.0"
const port = "8080"
var address = net.JoinHostPort(host, port)

func onRequest(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "<h1>Hello world!</h1>")
}

func main() {
	http.HandleFunc("/", onRequest)
	http.ListenAndServe(address, nil)
}
