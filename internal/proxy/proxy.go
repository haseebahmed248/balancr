// request forwarding
package proxy

import (
	"balancr/internal/logger"
	"balancr/internal/metrics"
	"balancr/internal/pool"
	"io"
	"net"
)

func SetupProxy(clientConn net.Conn, pools *pool.ServerPool, m *metrics.Metrics) {
	backend, err := pools.GetNext()
	if err != nil {
		logger.Log("No Backend Is Alive", "ERROR")
		m.IncrementError()
		clientConn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: 19\r\n\r\nNo backend is alive"))
		return
	}
	conn, err := net.Dial("tcp", backend)
	if err != nil {
		logger.Log(err.Error(), "ERROR")
		m.IncrementError()
		return
	}
	defer conn.Close()
	go io.Copy(conn, clientConn)
	logger.Log("Request forwarded to "+backend, "INFO")
	m.IncrementRequest(backend)
	io.Copy(clientConn, conn)

}
