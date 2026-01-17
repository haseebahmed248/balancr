package main

import (
	"balancr/internal/pool"
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

	backends := [...]string{"localhost:9999", "localhost:9998", "localhost:9997"}
	pool := pool.GetServerPool(backends[:])

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go proxy.SetupProxy(conn, pool)
	}
}

func main() {
	startServer()
}
