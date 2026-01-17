// request forwarding
package proxy

import (
	"balancr/internal/pool"
	"io"
	"log"
	"net"
)

func SetupProxy(clientConn net.Conn, pools *pool.ServerPool) {
	conn, err := net.Dial("tcp", pools.GetNext())
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()
	go io.Copy(conn, clientConn)
	io.Copy(clientConn, conn)

}
