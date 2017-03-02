package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/hdhauk/ShitChat/msg"
)

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

var users = threadSafeUsers{
	list: make(map[string]user),
	mu:   sync.Mutex{},
}

var newUserCh = make(chan user)
var broadcastChatMsg = make(chan chatMsg)

func main() {
	go incommingConnListenAndAccept(mux, "7000")

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

func mux(c net.Conn) {
	incommingReq := make(chan msg.ClientReq)
	outgoingResp := make(chan msg.ServerResp)
	closeConnCh := make(chan struct{})
	go rx(c, incommingReq)
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
			}
		}
	}
}

func handleLogout(username string, closeConnCh chan struct{}) {
	users.Remove(username)
	closeConnCh <- struct{}{}
}

func handleNames(out chan msg.ServerResp) {
	users.mu.Lock()
	usernames := users.list
	users.mu.Unlock()

	namesString := ""
	for k := range usernames {
		namesString = namesString + k + "\n"
	}
	out <- msg.ServerResp{Resp: "names", Content: namesString}
}

func handleMsg(message, username string, respCh chan msg.ServerResp) {
	broadcastChatMsg <- chatMsg{username: username, message: message}
}

func handleLogin(username string, respCh chan msg.ServerResp, closeConnCh chan struct{}) {
	// Check username validity
	if err := validate(username); err != nil {
		respCh <- msg.ServerResp{
			TimeStamp: time.Now().String(),
			Sender:    "server",
			Resp:      "Info",
			Content:   "Invalid username",
		}
		closeConnCh <- struct{}{}
		return
	}

	// Add user
	newUser := user{username, respCh}
	if err := users.Add(newUser); err != nil {
		respCh <- msg.ServerResp{
			TimeStamp: time.Now().String(),
			Sender:    "server",
			Resp:      "Error",
			Content:   err.Error(),
		}
		closeConnCh <- struct{}{}
		return
	}

	// Respond to user
	respCh <- msg.ServerResp{
		TimeStamp: time.Now().String(),
		Sender:    "server",
		Resp:      "Info",
		Content:   "Login successful",
	}

	// Send chat history
	hist := chatHistory.Dump()
	strHist := []string{}
	for _, v := range hist {
		strHist = append(strHist, fmt.Sprintf("[%v] %s", v.username, v.message))
	}
	respCh <- msg.ServerResp{
		TimeStamp: time.Now().String(),
		Sender:    "server",
		Resp:      "History",
		Content:   strHist,
	}
	log.Printf("[INFO] Chat history sent to %s\n", username)
}

func rx(c net.Conn, out chan msg.ClientReq) {
	emptyReq := msg.ClientReq{}
	for {
		// Decode incomming msg
		var req msg.ClientReq
		json.NewDecoder(c).Decode(&req)
		if req == emptyReq {
			continue
		}
		out <- req
	}
}

func tx(c net.Conn, out chan msg.ServerResp, closeConnCh chan struct{}) {
	for {
		select {
		case resp := <-out:
			// fmt.Printf("Sending %+v\n", resp)
			json.NewEncoder(c).Encode(&resp)
		case <-closeConnCh:
			log.Printf("[INFO] Closing connection to %+v\n", c.RemoteAddr())
			c.Close()
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
