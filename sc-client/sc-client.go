package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	"github.com/dimiro1/banner"
	"github.com/hdhauk/ShitChat/msg"
)

func rx(conn net.Conn) {
	emptyResp := msg.ServerResp{}
	for {
		var resp msg.ServerResp
		json.NewDecoder(conn).Decode(&resp)
		if resp == emptyResp {
			continue
		}
		switch strings.ToLower(resp.Resp) {
		case "names":
			fmt.Print(resp.Content)
		case "history":
			fmt.Println("=== Previous messages ===")
			if reflect.TypeOf(resp.Content).Kind() != reflect.Slice {
				fmt.Println("Received bad history")
			} else {
				msgs := resp.Content.([]interface{})
				for _, v := range msgs {
					fmt.Println(v)
				}
			}
			fmt.Println("=== Current messages ===")
		case "help":
			fmt.Println(resp.Content)
		case "message":
			fmt.Printf("%s [%s] %s: %s\n", resp.TimeStamp[:19], resp.Resp, resp.Sender, resp.Content)
		case "info":
			fmt.Printf("[INFO] %s\n", resp.Content)
		case "error":
			fmt.Printf("[ERROR] %s\n", resp.Content)
		default:
			fmt.Printf("[WARN] Unrecognized response: %+v\n", resp)
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
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(bannerTxt))
	var serverAddr = "localhost:7000"
	flag.StringVar(&serverAddr, "server", serverAddr, "ip:port to the server")
	flag.Parse()
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to server %s: %s\n", serverAddr, err.Error())
	}

	fmt.Println("=== ShitChat Client ===")
	fmt.Printf("Connected to server at %v\n", conn.RemoteAddr())

	outbox := make(chan msg.ClientReq)
	go rx(conn)
	go tx(conn, outbox)

	// Log on
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimRight(username, "\n")
	json.NewEncoder(conn).Encode(msg.ClientReq{Request: "login", Content: username})

	// Caputure Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		outbox <- msg.ClientReq{Request: "logout"}
		time.Sleep(200 * time.Millisecond)
		os.Exit(1)
	}()

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

var bannerTxt = `
{{ .AnsiColor.Cyan }}███████╗██╗  ██╗██╗████████╗     ██████╗██╗  ██╗ █████╗ ████████╗
██╔════╝██║  ██║██║╚══██╔══╝    ██╔════╝██║  ██║██╔══██╗╚══██╔══╝
███████╗███████║██║   ██║       ██║     ███████║███████║   ██║
╚════██║██╔══██║██║   ██║       ██║     ██╔══██║██╔══██║   ██║
███████║██║  ██║██║   ██║       ╚██████╗██║  ██║██║  ██║   ██║
╚══════╝╚═╝  ╚═╝╚═╝   ╚═╝        ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝
================================================================={{ .AnsiColor.Default }}
`
