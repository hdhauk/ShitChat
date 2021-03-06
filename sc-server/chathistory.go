package main

import (
	"fmt"
	"sync"
	"time"
)

type chatMsg struct {
	username string
	message  string
}

type threadSafeHistory struct {
	messages []chatMsg
	mu       sync.Mutex
}

func (tsh *threadSafeHistory) Add(cm chatMsg) {
	tsh.mu.Lock()
	// Add timestamps to messages
	cm.message = fmt.Sprintf("%s > %s", time.Now().String()[:19], cm.message)
	tsh.messages = append(tsh.messages, cm)
	tsh.mu.Unlock()
}

func (tsh *threadSafeHistory) Dump() []chatMsg {
	tsh.mu.Lock()
	defer tsh.mu.Unlock()
	return tsh.messages
}
