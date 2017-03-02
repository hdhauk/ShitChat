package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/hdhauk/ShitChat/msg"
)

type user struct {
	username string
	respCh   chan msg.ServerResp
}

type threadSafeUsers struct {
	list map[string]user
	mu   sync.Mutex
}

func (tsu *threadSafeUsers) Add(newUser user) error {
	tsu.mu.Lock()
	defer tsu.mu.Unlock()

	if u, ok := tsu.list[newUser.username]; ok {
		return fmt.Errorf("username %s already taken", u.username)
	}
	tsu.list[newUser.username] = newUser
	log.Printf("[INFO] Successfully added %s\n", newUser.username)
	return nil
}

func (tsu *threadSafeUsers) Remove(username string) {
	tsu.mu.Lock()
	defer tsu.mu.Unlock()
	delete(tsu.list, username)
	log.Printf("[INFO] Successfully removed %s\n", username)
}

func (tsu *threadSafeUsers) DumpAllUsernames() []string {
	tsu.mu.Lock()
	defer tsu.mu.Unlock()
	ret := []string{}
	for k := range tsu.list {
		ret = append(ret, k)
	}
	return ret
}

func (tsu *threadSafeUsers) DumpAllUsers() []user {
	tsu.mu.Lock()
	defer tsu.mu.Unlock()
	ret := []user{}
	for _, u := range tsu.list {
		ret = append(ret, u)
	}
	return ret
}
