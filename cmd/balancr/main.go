package main

import (
	"balancr/internal/config"
	"balancr/internal/logger"
	"balancr/internal/metrics"
	"balancr/internal/pool"
	"balancr/internal/proxy"
	"fmt"
	"net"
	"time"
)

func startServer(m *metrics.Metrics) {
	listener, err := net.Listen("tcp", "localhost:7000")
	if err != nil {
		logger.Log(err.Error(), "ERROR")
	}
	logger.Log("Listening to port 7000", "INFO")

	backends := config.GetBackends()

	data := make([]*pool.Backend, len(backends))
	for i, backend := range backends {
		data[i] = &pool.Backend{
			URL:    backend.URL,
			Weight: backend.Weight,
			Quota:  backend.Weight,
			Alive:  true,
		}
		pool.IsAlive(data[i])
	}
	serverPool := pool.GetServerPool(data)

	go serverPool.HealthCheck(10 * time.Second)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Log(err.Error(), "ERROR")
		}
		go proxy.SetupProxy(conn, serverPool, m)
	}
}

func serverMetrics(m *metrics.Metrics) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		logger.Log(fmt.Sprint(m.GetMetrics()), "INFO")
	}
}

func main() {
	var m = &metrics.Metrics{}
	go serverMetrics(m)
	startServer(m)

}
