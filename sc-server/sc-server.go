package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/hdhauk/ShitChat/msg"
)

// Global storage
var users = threadSafeUsers{list: make(map[string]user), mu: sync.Mutex{}}
var chatHistory = threadSafeHistory{}

// Internal server communication
var broadcastChatMsg = make(chan chatMsg)

func main() {
	listenPort := "7000"
	flag.StringVar(&listenPort, "port", listenPort, "Port for the server to listen on")
	flag.Parse()

	go incommingConnListenAndAccept(socketHandler, listenPort)

	// Keep track of all registered users
	for {
		select {
		// Broadcast chatmessages to all registered users
		case m := <-broadcastChatMsg:
			toSend := msg.ServerResp{
				TimeStamp: time.Now().String(),
				Sender:    m.username,
				Resp:      "Message",
				Content:   m.message,
			}
			// Broadcast to all users
			allUsers := users.DumpAllUsers()
			for _, u := range allUsers {
				u.respCh <- toSend
			}
			// Save to history
			chatHistory.Add(m)

		}
	}
}

func incommingConnListenAndAccept(handleConn func(c net.Conn), port string) {
	localAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConn(conn)
	}
}

func socketHandler(c net.Conn) {
	incommingReq := make(chan msg.ClientReq)
	outgoingResp := make(chan msg.ServerResp)
	closeConnCh := make(chan struct{})
	go rx(c, incommingReq, closeConnCh)
	go tx(c, outgoingResp, closeConnCh)

	var username string

	for {
		select {
		// Handle incomming requests
		case incomming := <-incommingReq:
			switch incomming.Request {
			case "login":
				username = incomming.Content.(string)
				go handleLogin(username, outgoingResp, closeConnCh)
			case "msg":
				go handleMsg(incomming.Content.(string), username, outgoingResp)
			case "names":
				handleNames(outgoingResp)
			case "logout":
				handleLogout(username, closeConnCh)
			case "help":
				handleHelp(outgoingResp)
			}
		}
	}
}

func rx(c net.Conn, out chan msg.ClientReq, closeConnCh chan struct{}) {
	emptyReq := msg.ClientReq{}
	remoteAddr := c.RemoteAddr().String()
	for {
		// Decode incomming msg
		var req msg.ClientReq
		err := json.NewDecoder(c).Decode(&req)
		if err != nil {
			log.Printf("[ERROR] Decoding request: %s\n", err.Error())
			close(closeConnCh)
			log.Printf("[INFO] Stopping rx for client %v\n", remoteAddr)
			return
		}
		if req == emptyReq {
			continue
		}
		out <- req
	}
}

func tx(c net.Conn, out chan msg.ServerResp, closeConnCh chan struct{}) {
	remoteAddr := c.RemoteAddr().String()
	for {
		select {
		case resp := <-out:
			// fmt.Printf("Sending %+v\n", resp)
			json.NewEncoder(c).Encode(&resp)
		case <-closeConnCh:
			log.Printf("[INFO] Closing connection to %+v\n", remoteAddr)
			c.Close()
			log.Printf("[INFO] Stopping tx for client %v\n", remoteAddr)
			return
		}
	}
}

func validate(username string) error {
	if !valid.IsAlphanumeric(username) || username == "" {
		return fmt.Errorf("invalid username")
	}
	return nil
}
