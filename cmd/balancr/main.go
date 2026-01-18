package main

import (
	"balancr/internal/pool"
	"balancr/internal/proxy"
	"log"
	"net"
	"time"
)

func startServer() {
	listener, err := net.Listen("tcp", "localhost:7000")
	if err != nil {
		log.Print(err)
	}
	log.Print("Listening to port 7000")

	backends := [...]string{"localhost:9999", "localhost:9998", "localhost:9997"}
	data := make([]*pool.Backend, len(backends))
	for i, backend := range backends {
		data[i] = &pool.Backend{
			URL:   backend,
			Alive: true,
		}
		pool.IsAlive(data[i])
	}
	serverPool := pool.GetServerPool(data)

	go serverPool.HealthCheck(10 * time.Second)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go proxy.SetupProxy(conn, serverPool)
	}
}

func main() {
	startServer()
}
