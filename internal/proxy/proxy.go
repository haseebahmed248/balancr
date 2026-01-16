// request forwarding
package proxy

import (
	"io"
	"log"
	"net"
)

func SetupProxy(clientConn net.Conn) {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	go io.Copy(conn, clientConn)
	io.Copy(clientConn, conn)
}
