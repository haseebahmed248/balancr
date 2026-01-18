// request forwarding
package proxy

import (
	"balancr/internal/pool"
	"io"
	"log"
	"net"
)

func SetupProxy(clientConn net.Conn, pools *pool.ServerPool) {
	backend, err := pools.GetNext()
	if err != nil {
		log.Print("No backend is running")
		clientConn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: 19\r\n\r\nNo backend is alive"))
		return
	}
	conn, err := net.Dial("tcp", backend)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	go io.Copy(conn, clientConn)
	io.Copy(clientConn, conn)

}
