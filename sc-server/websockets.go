package main

import (
	"log"
	"net/http"
	"strings"
)

type wsHub struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

func newHub() *wsHub {
	return &wsHub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

func (wsh *wsHub) run() {
	for {
		select {
		case client := <-wsh.register:
			wsh.clients[client] = true
		case client := <-wsh.unregister:
			if _, ok := wsh.clients[client]; ok {
				delete(wsh.clients, client)
				close(client.txCh)
			}
		case msg := <-wsh.broadcast:
			if strings.HasPrefix(string(msg), "[WebsocketUser]") {
				continue
			}
			for client := range wsh.clients {
				select {
				case client.txCh <- msg:
				default:
					close(client.txCh)
					delete(wsh.clients, client)
				}
			}
			broadcastChatMsg <- chatMsg{username: "WebsocketUser", message: string(msg)}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}
