package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"bitbucket.org/halvor_haukvik/ttm4100-go/msg"
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

func printContent(c net.Conn) {
	var msg msg.ClientReq
	json.NewDecoder(c).Decode(&msg)
	fmt.Printf("%+v", msg)
}

func main() {
	go incommingConnListenAndAccept(func(c net.Conn) { printContent(c) }, "7000")

	// Block forever
	select {}
}
