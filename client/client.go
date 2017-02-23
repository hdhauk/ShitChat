package main

import (
	"bytes"
	"encoding/json"
	"net"

	"bitbucket.org/halvor_haukvik/ttm4100-go/msg"
)

func sendHello() {
	conn, _ := net.Dial("tcp", "localhost:7000")
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msg.ClientReq{Request: "hello", Content: "test"})
	conn.Write(b.Bytes())
}

func main() {
	sendHello()
}
