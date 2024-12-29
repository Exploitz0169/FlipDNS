package main

import (
	"fmt"
	"net"

	"github.com/exploitz0169/flipdns/internal/udpserver"
)

func main() {

	addr := ":53"

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	fmt.Printf("Listening on %s\n", addr)
	udpserver.Run(conn)

}
