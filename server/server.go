package main

import (
	"fmt"
	"log"
	"net"
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

}

func main() {
	go incommingConnListenAndAccept(func(c net.Conn) { fmt.Println(c) }, "7000")

	// Block forever
	select {}
}
