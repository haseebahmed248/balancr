package main

import (
	"balancr/internal/proxy"
	"log"
	"net"
)

func startServer() {
	listener, err := net.Listen("tcp", "localhost:7000")
	if err != nil {
		log.Print(err)
	}
	log.Print("Listening to port 7000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go proxy.SetupProxy(conn)
	}
}

func main() {
	startServer()
}
