package main

import "sync"

type chatMsg struct {
	username string
	message  string
}

type history struct {
	messages []chatMsg
	mu       sync.Mutex
}

func (h *history) Add(cm chatMsg) {
	chatHistory.mu.Lock()
	chatHistory.messages = append(chatHistory.messages, cm)
	chatHistory.mu.Unlock()
}

func (h *history) Dump() []chatMsg {
	chatHistory.mu.Lock()
	defer chatHistory.mu.Unlock()
	return chatHistory.messages
}

var chatHistory = history{}
