package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hdhauk/ShitChat/msg"
)

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
