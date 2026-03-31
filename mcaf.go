package mcaf

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	localAddr := "0.0.0.0:25565"
	remoteAddr := "127.0.0.1:25560"

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}
	defer listener.Close()

	fmt.Printf("TCP Proxy listening on %s, forwarding to %s\n", localAddr, remoteAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, remoteAddr)
	}
}

func handleConnection(clientConn net.Conn, remoteAddr string) {
	defer clientConn.Close()
	clientIP := clientConn.RemoteAddr().String()

	log.Printf("%s", clientIP)

	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("Failed to connect to remote server %s: %v", remoteAddr, err)
		return
	}
	defer remoteConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		io.Copy(remoteConn, clientConn)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(clientConn, remoteConn)
		done <- struct{}{}
	}()

	<-done
}
