package main

import (
	"balancr/internal/config"
	"balancr/internal/logger"
	"balancr/internal/metrics"
	"balancr/internal/pool"
	"balancr/internal/proxy"
	"flag"
	"fmt"
	"net"
	"time"
)

func startServer(m *metrics.Metrics, port int, config_file string, health_interval int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		logger.Log(err.Error(), "ERROR")
	}
	logger.Log(fmt.Sprintf("Listening to port %d", port), "INFO")

	backends := config.GetBackends(config_file)

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

	go serverPool.HealthCheck(time.Duration(health_interval) * time.Second)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Log(err.Error(), "ERROR")
		}
		go proxy.SetupProxy(conn, serverPool, m)
	}
}

func serverMetrics(m *metrics.Metrics, metrics_interval int) {
	ticker := time.NewTicker(time.Duration(metrics_interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		logger.Log(fmt.Sprint(m.GetMetrics()), "INFO")
	}
}

func main() {

	port := flag.Int("port", 7000, "port on which proxy will run")
	config := flag.String("config", "config.yaml", "Your config.yaml file which holds the backends")
	health_interval := flag.Int("health-interval", 10, "Interval after which server health will be checked")
	metrics_interval := flag.Int("metrics-interval", 10, "Interval after which metrics will be displayed")
	flag.Parse()
	var m = &metrics.Metrics{}
	go serverMetrics(m, *metrics_interval)
	startServer(m, *port, *config, *health_interval)

}
