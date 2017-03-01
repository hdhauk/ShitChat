package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"bitbucket.org/halvor_haukvik/ttm4100-go/msg"
)

func rx(conn net.Conn) {
	emptyResp := msg.ServerResp{}
	for {
		var resp msg.ServerResp
		json.NewDecoder(conn).Decode(&resp)
		if resp == emptyResp {
			continue
		}
		switch resp.Resp {
		case "names":
			fmt.Print(resp.Content)
		case "History":
			fmt.Println("=== Previous messages ===")
			msgs := resp.Content.([]interface{})
			for _, v := range msgs {
				fmt.Println(v)
			}
			fmt.Println("=== Current messages ===")
		default:
			fmt.Printf("%s [%s] %s: %s\n", resp.TimeStamp, resp.Resp, resp.Sender, resp.Content)
		}
	}
}

func tx(conn net.Conn, outbox chan msg.ClientReq) {
	for {
		select {
		case tx := <-outbox:
			json.NewEncoder(conn).Encode(&tx)
		}
	}

}

func main() {
	conn, _ := net.Dial("tcp", "localhost:7000")

	outbox := make(chan msg.ClientReq)
	go rx(conn)
	go tx(conn, outbox)

	// Log on
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimRight(username, "\n")
	json.NewEncoder(conn).Encode(msg.ClientReq{Request: "login", Content: username})

loop:
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimRight(message, "\n")

		switch message {
		case "\\quit":
			outbox <- msg.ClientReq{Request: "logout"}
			time.Sleep(500 * time.Millisecond)
			break loop
		case "\\names":
			outbox <- msg.ClientReq{Request: "names"}
		case "\\help":
			outbox <- msg.ClientReq{Request: "help"}
		default:
			outbox <- msg.ClientReq{Request: "msg", Content: message}

		}
	}
}
